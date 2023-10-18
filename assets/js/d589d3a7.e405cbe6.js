"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7162],{3905:(e,t,n)=>{n.d(t,{Zo:()=>d,kt:()=>g});var a=n(7294);function l(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function r(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){l(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,a,l=function(e,t){if(null==e)return{};var n,a,l={},i=Object.keys(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||(l[n]=e[n]);return l}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(a=0;a<i.length;a++)n=i[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(l[n]=e[n])}return l}var s=a.createContext({}),c=function(e){var t=a.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):r(r({},t),e)),n},d=function(e){var t=c(e.components);return a.createElement(s.Provider,{value:t},e.children)},p="mdxType",u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},m=a.forwardRef((function(e,t){var n=e.components,l=e.mdxType,i=e.originalType,s=e.parentName,d=o(e,["components","mdxType","originalType","parentName"]),p=c(n),m=l,g=p["".concat(s,".").concat(m)]||p[m]||u[m]||i;return n?a.createElement(g,r(r({ref:t},d),{},{components:n})):a.createElement(g,r({ref:t},d))}));function g(e,t){var n=arguments,l=t&&t.mdxType;if("string"==typeof e||l){var i=n.length,r=new Array(i);r[0]=m;var o={};for(var s in t)hasOwnProperty.call(t,s)&&(o[s]=t[s]);o.originalType=e,o[p]="string"==typeof e?e:l,r[1]=o;for(var c=2;c<i;c++)r[c]=n[c];return a.createElement.apply(null,r)}return a.createElement.apply(null,n)}m.displayName="MDXCreateElement"},9390:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>r,default:()=>u,frontMatter:()=>i,metadata:()=>o,toc:()=>c});var a=n(7462),l=(n(7294),n(3905));const i={sidebar_position:1},r="Getting started",o={unversionedId:"getting-started",id:"getting-started",title:"Getting started",description:"Prerequisites:",source:"@site/docs/getting-started.md",sourceDirName:".",slug:"/getting-started",permalink:"/docs/getting-started",draft:!1,editUrl:"https://github.com/scaling-lightning/scaling-lightning/tree/main/docs/docs/getting-started.md",tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1},sidebar:"docsSidebar",next:{title:"Architectural overview",permalink:"/docs/architectural-overview"}},s={},c=[{value:"Prerequisites:",id:"prerequisites",level:2},{value:"Installation",id:"installation",level:2},{value:"Starting a Network",id:"starting-a-network",level:2},{value:"Example CLI Commands",id:"example-cli-commands",level:2},{value:"Run the above from code instead of CLI",id:"run-the-above-from-code-instead-of-cli",level:2},{value:"Helpful Kubernetes commands",id:"helpful-kubernetes-commands",level:2}],d={toc:c},p="wrapper";function u(e){let{components:t,...n}=e;return(0,l.kt)(p,(0,a.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"getting-started"},"Getting started"),(0,l.kt)("h2",{id:"prerequisites"},"Prerequisites:"),(0,l.kt)("ul",null,(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},"Kubernetes."),(0,l.kt)("ul",{parentName:"li"},(0,l.kt)("li",{parentName:"ul"},"If you are developing locally you can use Docker Desktop and enable\nKubernetes in the dashboard."),(0,l.kt)("li",{parentName:"ul"},"Alternatively minikube works as an alternative to Docker Desktop. Please use ",(0,l.kt)("inlineCode",{parentName:"li"},"minikube tunnel"),' to enable traefik to get an "external" ip which the library and cli requires to communicate in to the sidecar clients.'),(0,l.kt)("li",{parentName:"ul"},"SL has also been tested on Digital Ocean's hosted K8s cluster"),(0,l.kt)("li",{parentName:"ul"},"Please let us know if you have run SL on a different cluster distribution such as Kind, K3s K0s or any other cloud provider"))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},"Helm 3 and Helmfile."),(0,l.kt)("p",{parentName:"li"},"Mac OS"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre"},"brew install helm helmfile\n")),(0,l.kt)("p",{parentName:"li"},"Windows"),(0,l.kt)("pre",{parentName:"li"},(0,l.kt)("code",{parentName:"pre"},"scoop install helm helmfile\n")),(0,l.kt)("p",{parentName:"li"},"For Linux check your distros package manager but you may need to download the binaries for helm and helmfile.")),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},"Helm Diff:"),(0,l.kt)("p",{parentName:"li"},"  helm plugin install ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/databus23/helm-diff"},"https://github.com/databus23/helm-diff")),(0,l.kt)("blockquote",{parentName:"li"},(0,l.kt)("p",{parentName:"blockquote"},(0,l.kt)("strong",{parentName:"p"},(0,l.kt)("em",{parentName:"strong"},"NOTE:"))," On Windows the plugin install does not complete correctly and you need to download the binary manually from ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/databus23/helm-diff/releases"},"https://github.com/databus23/helm-diff/releases")," . Unzip the diff.exe file and put it in the ",(0,l.kt)("em",{parentName:"p"},"helm/plugins/helm-diff/bin")," folder (the ",(0,l.kt)("em",{parentName:"p"},"bin")," folder has to be created). You can find the folder by running ",(0,l.kt)("em",{parentName:"p"},'"helm env HELM_DATA_HOME"')))),(0,l.kt)("li",{parentName:"ul"},(0,l.kt)("p",{parentName:"li"},"Traefik:"),(0,l.kt)("p",{parentName:"li"},"  helm repo add traefik ",(0,l.kt)("a",{parentName:"p",href:"https://traefik.github.io/charts"},"https://traefik.github.io/charts"),"\nhelm repo update\nhelm install traefik traefik/traefik -n sl-traefik --create-namespace -f ",(0,l.kt)("a",{parentName:"p",href:"https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/charts/traefik-values.yml"},"https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/charts/traefik-values.yml")))),(0,l.kt)("h2",{id:"installation"},"Installation"),(0,l.kt)("p",null,"Download binary for your system from ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/scaling-lightning/scaling-lightning/releases"},"Releases")),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"# untar to get binary\ntar -xzf scaling-lightning-[version]-[os]-[architecture].tar.gz\n\n# Mac OS only - mark file as safe so it will run\nxattr -dr com.apple.quarantine scaling-lightning\n\n# run - should print CLI help\n./scaling-lightning\n")),(0,l.kt)("h2",{id:"starting-a-network"},"Starting a Network"),(0,l.kt)("p",null,"To spin up an example network with 2 cln nodes and 4 lnd nodes, run:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"# Download example helmfile which defines the nodes you want in your network.\nwget https://raw.githubusercontent.com/scaling-lightning/scaling-lightning/main/examples/helmfiles/public.yaml\n\n# Create and start the network. Scaling lightning will use your currently defined default k8s cluster\n# as specified in kubectl kubectl config get-contexts\n./scaling-lightning create -f public.yaml\n")),(0,l.kt)("p",null,"To destroy the network run:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"./scaling-lightning destroy\n")),(0,l.kt)("h2",{id:"example-cli-commands"},"Example CLI Commands"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"# list nodes on the network (names were taken from the helmfile)\n./scaling-lightning list\n\n# get wallet balance of node named bitcoind\n./scaling-lightning walletbalance -n bitcoind\n\n# get wallet balance of node named lnd2\n./scaling-lightning walletbalance -n lnd2\n\n# send on-chain 1 million satoshis from bitcoind to cln1\n./scaling-lightning send -f bitcoind -t cln1 -a 1000000\n\n# get the pubkey of a node named lnd1\n./scaling-lightning pubkey -n lnd1\n\n# peer lnd1 and cln1 from lnd1\n./scaling-lightning connectpeer -f lnd1 -t cln1\n\n# open channel between cln1 and lnd1 with a local balance on cln1 of 70k satoshis\n./scaling-lightning openchannel -f cln1 -t lnd1 -a 70000\n\n# have bitcoind generate some blocks and pay itself the block reward\n./scaling-lightning generate -n bitcoind\n")),(0,l.kt)("h2",{id:"run-the-above-from-code-instead-of-cli"},"Run the above from code instead of CLI"),(0,l.kt)("p",null,"See ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/scaling-lightning/scaling-lightning/blob/main/examples/go/example_test.go"},"examples/go/example_test.go"),". This test takes around 3 minutes to pass on an M1 Macbook Pro so you may need to adjust your test runner's default timeout."),(0,l.kt)("p",null,"Example go test command with extra timeout:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"go test -run ^TestMain$ github.com/scaling-lightning/scaling-lightning/examples/go -count=1 -v -timeout=15m\n")),(0,l.kt)("h2",{id:"helpful-kubernetes-commands"},"Helpful Kubernetes commands"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"# list pods\nkubectl -n sl get pods\n\n# describe cln1 pod in more detail\nkubectl -n sl describe pod cln1-0\n\n# view logs of lnd1 node\nkubectl -n sl logs -f lnd1-0\n\n# view logs of a crashed bitcoind pod\nkubectl -n sl logs -previous bitcoind-0\n\n# view logs of lnd1's scaling lightning sidecar client (it handles our api requests and forwards them to the node)\nkubectl -n sl logs -f -c lnd-client lnd1-0\n\n# same for cln and bitcoind\nkubectl -n sl logs -f -c cln-client cln1-0\nkubectl -n sl logs -f -c bitcoind-client bitcoind-0\n\n# get shell into lnd1\nkubectl -n sl exec -it lnd1-0 -- bash\n\n# view loadbalancer public ip from traefik\nkubectl -n sl-traefik get services\n\n# destroy all scaling lightning nodes\nkubectl delete namespace sl\n\n# uninstall traefik\nkubectl delete namespace sl-traefik\n\n# uninstall traefik alternative\nhelm uninstall traefik -n sl-traefik\n")),(0,l.kt)("p",null,"Note that the above commands assume you are using the default kubeconfig and context. You would need to add ",(0,l.kt)("inlineCode",{parentName:"p"},"--kubeconfig path/to/file.yml")," or ",(0,l.kt)("inlineCode",{parentName:"p"},"--context mycluster")," to all of the above commands if you wanted to look at a different cluster."))}u.isMDXComponent=!0}}]);