package builder

import (
	"encoding/json"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func TaskForMdi(srcCheat string, destCheat string, res string, gofile string) {
	initRemixResourceTemplate(res, gofile)
	writeRemixCheatSheet(res, destCheat)
	_ = srcCheat // kept for call-site shape; cheat sheet is generated from paths
}

func initRemixResourceTemplate(src string, dest string) {
	// https://github.com/Remix-Design/RemixIcon fonts/remixicon.glyph.json (or slim paths.json)
	fileRaw, err := os.ReadFile(filepath.Clean(src))
	if err != nil {
		fmt.Println("读取文件出错", src, err)
		return
	}

	icons, err := parseRemixPaths(fileRaw)
	if err != nil {
		fmt.Println("解析 RemixIcon 出错", err)
		return
	}

	raw, _ := json.MarshalIndent(icons, "", " ")
	goFile := "package mdi\nvar iconMap = map[string]string" + string(raw)
	goFile = strings.Replace(goFile, "\"\n}", "\",\n}", 1)
	content, err := format.Source([]byte(goFile))
	if err != nil {
		fmt.Println("格式化 icons.go 出错", err)
		return
	}

	if err := os.WriteFile(dest, content, os.ModePerm); err != nil {
		fmt.Println("保存文件出错", err)
	} else {
		fmt.Println("保存 RemixIcon 资源文件完毕", dest, "icons:", len(icons))
	}
}

func parseRemixPaths(raw []byte) (map[string]string, error) {
	// slim form: {"home-line":"M..."}
	var slim map[string]string
	if err := json.Unmarshal(raw, &slim); err == nil {
		if len(slim) > 0 {
			// reject if values look like nested objects stringified — treat as ok
			out := make(map[string]string, len(slim))
			for k, v := range slim {
				if v == "" {
					continue
				}
				out[strings.ToLower(k)] = v
			}
			if len(out) > 0 {
				// if this was actually glyph.json, values wouldn't be plain paths; check one value
				for _, v := range out {
					if strings.HasPrefix(strings.TrimSpace(v), "{") {
						out = nil
						break
					}
					break
				}
				if out != nil {
					return out, nil
				}
			}
		}
	}

	// full glyph form: {"home-line":{"path":["M..."],...}}
	var glyph map[string]struct {
		Path []string `json:"path"`
	}
	if err := json.Unmarshal(raw, &glyph); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(glyph))
	for k, v := range glyph {
		if len(v.Path) == 0 {
			continue
		}
		out[strings.ToLower(k)] = strings.Join(v.Path, " ")
	}
	return out, nil
}

func writeRemixCheatSheet(res string, dest string) {
	fileRaw, err := os.ReadFile(filepath.Clean(res))
	if err != nil {
		log.Println("生成图标页失败: 读取资源", err)
		return
	}
	icons, err := parseRemixPaths(fileRaw)
	if err != nil {
		log.Println("生成图标页失败: 解析", err)
		return
	}

	names := make([]string, 0, len(icons))
	for name := range icons {
		names = append(names, name)
	}
	sort.Strings(names)

	// names-only index + lazy path lookup keeps the embed small; paths live in icons.go already via GetIconByName flows.
	// For the browser page we still need paths — put them once as JSON (smaller than 3k HTML nodes).
	type pair struct{ N, P string }
	list := make([]pair, 0, len(names))
	for _, name := range names {
		list = append(list, pair{N: name, P: icons[name]})
	}
	payload, _ := json.Marshal(list)

	var b strings.Builder
	b.WriteString(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>Remix Icon</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{font-family:system-ui,sans-serif;background:#fafafa;color:#222;padding:16px}
h1{font-size:20px;margin-bottom:8px}
p{color:#666;margin-bottom:16px;font-size:14px}
input{width:100%;max-width:420px;padding:8px 12px;border:1px solid #ddd;border-radius:8px;margin-bottom:16px;font-size:14px}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(120px,1fr));gap:8px}
.item{display:flex;flex-direction:column;align-items:center;gap:6px;padding:12px 8px;background:#fff;border:1px solid #eee;border-radius:8px;cursor:pointer}
.item:hover{border-color:#006aff}
.item svg{width:28px;height:28px;fill:currentColor}
.item span{font-size:11px;word-break:break-all;text-align:center;color:#555}
.hidden{display:none}
</style>
</head>
<body>
<h1>Remix Icon</h1>
<p>Click name to copy. Use as bookmark/app icon field. Source: <a href="https://remixicon.com">remixicon.com</a></p>
<input id="q" type="search" placeholder="filter icons…" autofocus/>
<div class="grid" id="grid"></div>
<script type="application/json" id="data">`)
	b.Write(payload)
	b.WriteString(`</script>
<script>
const data=JSON.parse(document.getElementById('data').textContent);
const grid=document.getElementById('grid'),q=document.getElementById('q');
const frag=document.createDocumentFragment();
for(const {N,P} of data){
  const el=document.createElement('div');el.className='item';el.dataset.name=N;el.title=N;
  el.innerHTML='<svg viewBox="0 0 24 24"><path d="'+P+'"></path></svg><span></span>';
  el.lastChild.textContent=N;frag.appendChild(el);
}
grid.appendChild(frag);
q.addEventListener('input',()=>{const v=q.value.trim().toLowerCase();
for(const el of grid.children)el.classList.toggle('hidden',v&&!el.dataset.name.includes(v));});
grid.addEventListener('click',async e=>{const el=e.target.closest('.item');if(!el)return;
try{await navigator.clipboard.writeText(el.dataset.name);el.style.outline='2px solid #006aff';setTimeout(()=>el.style.outline='',400);}catch(_){}});
</script>
</body>
</html>
`)

	_PrepareDirectory(dest)
	out := filepath.Join(dest, "index.html")
	if err := os.WriteFile(out, []byte(b.String()), 0644); err != nil {
		log.Println("写入图标页失败", err)
	} else {
		fmt.Println("保存 RemixIcon 浏览页完毕", out)
	}
}
