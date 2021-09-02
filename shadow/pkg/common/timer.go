package common

import (
	"log"
	"reflect"
	"sync"
	"time"
)

// CallbackFunc is function which will be invoked periodically
type CallbackFunc func(
	// context used to pass value which the caller  of timer  want to delivery
	context map[string]string,
	// executeContext used to save value which was passed between of a
	// CallbackFunc in each timer call
	executeContext map[string]interface{},
	// interval represented a time duration between each timer call.
	interval time.Duration,
) bool

type timerCallback struct {
	group    string
	context  map[string]string
	interval time.Duration
	callback CallbackFunc
	stop     chan bool
}

//Start the timerCallback, execute registered func
func (tcb *timerCallback) Start(registry *CallbackRegistry) {
	go func() {
		executeContext := make(map[string]interface{})
		timer := time.NewTicker(tcb.interval)

		// we need stop ticker explicitly
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				end := !tcb.callback(tcb.context, executeContext, tcb.interval)
				if end {
					registry.Remove(tcb.group, tcb.context, tcb.callback)
				}
			case <-tcb.stop:
				log.Printf("Timer [context: %+v group: %s interval: %+v] was removed ",
					tcb.context, tcb.group, tcb.interval)
				return
			}
		}
	}()
}

//Stop the timerCallback
func (tcb *timerCallback) Stop() {
	tcb.stop <- true
}

//Close the timerCallback
func (tcb *timerCallback) Close() {
	close(tcb.stop)
}

func newTimerCallback(group string, context map[string]string,
	interval time.Duration, callback CallbackFunc) *timerCallback {
	return &timerCallback{
		group:    group,
		context:  context,
		interval: interval,
		callback: callback,
		stop:     make(chan bool, 1),
	}
}

// CallbackRegistry is a registry which store timer for calling CallbackFunc
// periodically
type CallbackRegistry struct {
	sync.RWMutex
	// group by groupName (stack id, component name ro resource id)
	callbacks map[string][]*timerCallback
}

// Add used to add new Callback to the registry
func (r *CallbackRegistry) Add(group string, context map[string]string,
	interval time.Duration, callback CallbackFunc) bool {
	r.Lock()

	if interval <= 0 {
		log.Printf("invalid interval %d of timer registered callback", interval)
		return false
	}
	callbackList, existing := r.callbacks[group]
	if existing {
		for _, tcb := range callbackList {
			if tcb.group == group && reflect.DeepEqual(tcb.context, context) &&
				CompareFuncs(tcb.callback, callback) {
				// prevents create duplicated timer on same the group and callback
				log.Printf("time registry %+v duplicated, new timer was discard.", tcb)
				r.Unlock()
				return false
			}
		}
	}

	tcb := newTimerCallback(group, context, interval, callback)
	callbackList = append(callbackList, tcb)
	r.callbacks[group] = callbackList

	r.Unlock()

	tcb.Start(r)

	return true
}

// Remove used to remove a CallbackFunc from registry
func (r *CallbackRegistry) Remove(group string, context map[string]string,
	callback CallbackFunc) *CallbackRegistry {
	r.Lock()
	defer r.Unlock()

	callbackList, existing := r.callbacks[group]
	if !existing {
		return nil
	}

	removeIndices := make([]int, 0)
	for idx, tcb := range callbackList {
		if reflect.DeepEqual(tcb.context, context) &&
			tcb.group == group &&
			CompareFuncs(tcb.callback, callback) {
			removeIndices = append(removeIndices, idx)
			log.Printf("context: %+v group %s callback will be removed", context, group)
		} else {
			removeIndices = append(removeIndices, -1)
		}
	}

	if size, newList := r.removeTimerCallbackByIndices(callbackList, removeIndices); size == 0 {
		delete(r.callbacks, group)
	} else {
		r.callbacks[group] = newList
	}
	return r
}

// RemoveByGroup remove call by group
func (r *CallbackRegistry) RemoveByGroup(group string) *CallbackRegistry {
	r.Lock()
	defer r.Unlock()

	callbackList, existing := r.callbacks[group]
	if !existing {
		return nil
	}

	removeIndices := make([]int, 0)
	for idx, tcb := range callbackList {
		if tcb.group == group {
			removeIndices = append(removeIndices, idx)
			log.Printf("group %s callback will be removed", group)
		} else {
			removeIndices = append(removeIndices, -1)
		}
	}

	if size, newList := r.removeTimerCallbackByIndices(callbackList, removeIndices); size == 0 {
		delete(r.callbacks, group)
	} else {
		r.callbacks[group] = newList
	}
	return r
}
func (r *CallbackRegistry) removeTimerCallbackByIndices(
	callbackList []*timerCallback, indices []int) (int, []*timerCallback) {
	backupCallbackList := make([]*timerCallback, 0)
	for idx, callback := range callbackList {
		if indices[idx] == -1 {
			backupCallbackList = append(backupCallbackList, callback)
		} else {
			callback.Close()
		}
	}
	return len(backupCallbackList), backupCallbackList
}

// RemoveAll used to remove a batch of CallbackFunc
func (r *CallbackRegistry) RemoveAll(callback CallbackFunc) *CallbackRegistry {
	r.Lock()
	defer r.Unlock()

	removeIndices := make(map[string][]int, 0)

	for group, callbackList := range r.callbacks {
		var indices []int
		for idx, tcb := range callbackList {
			if CompareFuncs(tcb.callback, callback) {
				indices = append(indices, idx)
			} else {
				indices = append(indices, -1)
			}
		}
		removeIndices[group] = indices
	}

	for group, rmIndices := range removeIndices {
		if len(rmIndices) > 0 {
			if size, newList := r.removeTimerCallbackByIndices(r.callbacks[group],
				rmIndices); size == 0 {
				delete(r.callbacks, group)
			} else {
				r.callbacks[group] = newList
			}
		}
	}

	return r
}

// Stop used to stop all timer in registry
func (r *CallbackRegistry) Stop() {
	r.Lock()
	defer r.Unlock()

	for _, callbackList := range r.callbacks {
		for _, tcb := range callbackList {
			tcb.Stop()
		}
	}
}

// Close used to close all timer in registry
func (r *CallbackRegistry) Close() {
	r.Lock()
	defer r.Unlock()

	for _, callbackList := range r.callbacks {
		for _, tcb := range callbackList {
			tcb.Close()
		}
	}
}

// ListGroup list group name of callbacks
func (r *CallbackRegistry) ListGroup() []string {
	var result []string
	for k := range r.callbacks {
		result = append(result, k)
	}
	return result
}

// NewCallbackRegistry used to create TimerRegistry
func NewCallbackRegistry() *CallbackRegistry {
	r := new(CallbackRegistry)
	r.callbacks = make(map[string][]*timerCallback)
	return r
}

// CompareFuncs is used to identify whether two functions is equal
func CompareFuncs(func1 interface{}, func2 interface{}) bool {
	// FIXME: comparing two funcs depend on undefined behavior, but it works.
	// Maybe we can find another solution to compare funcs in the future
	sf1 := reflect.ValueOf(func1)
	sf2 := reflect.ValueOf(func2)
	return sf1.Pointer() == sf2.Pointer()
}
