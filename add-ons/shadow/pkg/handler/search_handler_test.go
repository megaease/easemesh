package handler

import (
	"reflect"
	"sync"
	"testing"

	shadowfake "github.com/megaease/easemesh/mesh-shadow/pkg/handler/fake"
	"github.com/megaease/easemesh/mesh-shadow/pkg/object"
	appsV1 "k8s.io/api/apps/v1"
)

func Test_isShadowDeployment(t *testing.T) {
	deployment1 := shadowfake.NewSourceDeployment()
	deployment2 := shadowfake.NewShadowDeployment()

	type args struct {
		spec appsV1.DeploymentSpec
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test1",
			args: args{
				spec: deployment1.Spec,
			},
			want: false,
		},
		{
			name: "test2",
			args: args{
				spec: deployment2.Spec,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isShadowDeployment(tt.args.spec); got != tt.want {
				t.Errorf("isShadowDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShadowServiceDeploySearcher_Search(t *testing.T) {
	searchChan := make(chan interface{})
	defer close(searchChan)

	searcher := &ShadowServiceDeploySearcher{
		KubeClient: prepareClientForTest(),
		ResultChan: searchChan,
	}

	sourceDeployment := shadowfake.NewSourceDeployment()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case obj := <-searcher.ResultChan:
				if !reflect.DeepEqual(obj.(ShadowServiceBlock).deployObj, *sourceDeployment) {
					t.Errorf("Search Deployment Error, Searcher.Search() = %v, \n want %v", obj, sourceDeployment)
				}
				return
			}
		}
	}()

	shadowService := shadowfake.NewShadowService()
	objs := []object.ShadowService{shadowService}
	searcher.Search(objs)
	wg.Wait()
}
