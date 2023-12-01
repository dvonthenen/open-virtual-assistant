module github.com/dvonthenen/open-virtual-assistant

go 1.18

require (
	cloud.google.com/go/speech v1.19.0
	cloud.google.com/go/texttospeech v1.6.0
	github.com/antzucaro/matchr v0.0.0-20210222213004-b04723ef80f0
	github.com/deepgram/deepgram-go-sdk v1.0.0
	github.com/dvonthenen/chat-gpeasy v0.2.2
	github.com/faiface/beep v1.1.0
	github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9
	google.golang.org/genproto v0.0.0-20230530153820-e85fd2cbaebc
	google.golang.org/grpc v1.55.0
	k8s.io/klog/v2 v2.110.1
)

require (
	cloud.google.com/go v0.110.2 // indirect
	cloud.google.com/go/compute v1.19.3 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/longrunning v0.5.0 // indirect
	github.com/dvonthenen/websocket v1.5.1-dyv.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.3 // indirect
	github.com/hajimehoshi/oto v0.7.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sashabaranov/go-openai v1.7.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.9.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/image v0.0.0-20190802002840-cff245a6509b // indirect
	golang.org/x/mobile v0.0.0-20201217150744-e6ae53a27f4f // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

// replace github.com/deepgram-devs/deepgram-go-sdk => ../../deepgram-devs/deepgram-go-sdk
