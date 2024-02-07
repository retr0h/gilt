"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[915],{1855:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>l,contentTitle:()=>s,default:()=>g,frontMatter:()=>o,metadata:()=>a,toc:()=>c});var i=t(5893),r=t(1151);const o={sidebar_position:4},s="Usage",a={id:"usage",title:"Usage",description:"CLI",source:"@site/versioned_docs/version-2.1.1/usage.md",sourceDirName:".",slug:"/usage",permalink:"/gilt/2.1.1/usage",draft:!1,unlisted:!1,tags:[],version:"2.1.1",sidebarPosition:4,frontMatter:{sidebar_position:4},sidebar:"docsSidebar",previous:{title:"Configuration",permalink:"/gilt/2.1.1/configuration"},next:{title:"Testing",permalink:"/gilt/2.1.1/testing"}},l={},c=[{value:"CLI",id:"cli",level:2},{value:"Init Configuration",id:"init-configuration",level:3},{value:"Overlay Repository",id:"overlay-repository",level:3},{value:"Debug",id:"debug",level:3},{value:"Package",id:"package",level:2},{value:"Overlay Repository",id:"overlay-repository-1",level:3}];function d(e){const n={code:"code",h1:"h1",h2:"h2",h3:"h3",p:"p",pre:"pre",...(0,r.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"usage",children:"Usage"}),"\n",(0,i.jsx)(n.h2,{id:"cli",children:"CLI"}),"\n",(0,i.jsx)(n.h3,{id:"init-configuration",children:"Init Configuration"}),"\n",(0,i.jsx)(n.p,{children:"Initializes config file in the shell's current working directory:"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"gilt init\n"})}),"\n",(0,i.jsx)(n.h3,{id:"overlay-repository",children:"Overlay Repository"}),"\n",(0,i.jsx)(n.p,{children:"Overlay a remote repository into the destination provided."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"gilt overlay\n"})}),"\n",(0,i.jsx)(n.h3,{id:"debug",children:"Debug"}),"\n",(0,i.jsx)(n.p,{children:"Display the git commands being executed."}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-bash",children:"gilt --debug overlay\n"})}),"\n",(0,i.jsx)(n.h2,{id:"package",children:"Package"}),"\n",(0,i.jsx)(n.h3,{id:"overlay-repository-1",children:"Overlay Repository"}),"\n",(0,i.jsxs)(n.p,{children:["See example client in ",(0,i.jsx)(n.code,{children:"examples/go-client/"}),"."]}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-go",children:'func main() {\n\tdebug := true\n\tlogger := getLogger(debug)\n\n\tc := config.Repositories{\n\t\tDebug:   debug,\n\t\tGiltDir: "~/.gilt",\n\t\tRepositories: []config.Repository{\n\t\t\t{\n\t\t\t\tGit:     "https://github.com/retr0h/ansible-etcd.git",\n\t\t\t\tVersion: "77a95b7",\n\t\t\t\tDstDir:  "../tmp/retr0h.ansible-etcd",\n\t\t\t},\n\t\t},\n\t}\n\n\tvar r repositoriesManager = repositories.New(c, logger)\n\tr.Overlay()\n}\n'})})]})}function g(e={}){const{wrapper:n}={...(0,r.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},1151:(e,n,t)=>{t.d(n,{Z:()=>a,a:()=>s});var i=t(7294);const r={},o=i.createContext(r);function s(e){const n=i.useContext(o);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:s(e.components),i.createElement(o.Provider,{value:n},e.children)}}}]);