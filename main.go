package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/onix-engine/onix/html/markup"
	"github.com/pencode-labs/pencode/peb"
)

var menuHTML = func() string {
	var b markup.Builder
	b.TextElement("button", "打开表单弹窗", "data-ref", "dialog", "type", "button", "role", "menuitem")
	b.TextElement("button", "打开自由定位浮窗", "data-ref", "popup", "type", "button", "role", "menuitem")
	b.TextElement("button", "打开工作台 Panel", "data-ref", "panel", "type", "button", "role", "menuitem")
	b.TextElement("button", "显示通知与自定义操作", "data-ref", "notification", "type", "button", "role", "menuitem")
	b.TextElement("button", "显示自动关闭通知", "data-ref", "auto-close", "type", "button", "role", "menuitem")
	return b.String()
}()
var formHTML = func() string {
	var b markup.Builder
	b.Open("form", "data-ref", "form")
	b.TextElement("label", "名称", "for", "pencode-extension-demo-name")
	b.Element("input", "id", "pencode-extension-demo-name", "data-ref", "name", "name", "name", "placeholder", "输入名称")
	b.TextElement("button", "在 Panel 中显示", "type", "submit")
	b.End()
	return b.String()
}()
var temporaryHTML = func() string {
	var b markup.Builder
	b.TextElement("p", "这条通知将在 5 秒后自动关闭，也可以点击右上角关闭。")
	return b.String()
}()

var notificationHTML = func() string {
	var b markup.Builder
	b.TextElement("p", "扩展可以自定义通知内容和操作。此通知会保留到手动关闭。")
	b.Open("footer")
	b.TextElement("button", "查看结果", "data-ref", "result", "type", "button")
	b.TextElement("button", "忽略", "data-ref", "dismiss", "type", "button")
	b.End()
	return b.String()
}()

var popupHTML = func() string {
	var b markup.Builder
	b.TextElement("p", "可像扩展市场一样显示丰富内容，同时继续操作工作区。")
	b.TextElement("label", "搜索", "for", "pencode-extension-demo-query")
	b.Element("input", "id", "pencode-extension-demo-query", "data-ref", "query", "placeholder", "移动浮窗后保留输入")
	b.Open("label")
	b.Element("input", "data-ref", "outside", "type", "checkbox")
	b.Text("点击外部关闭")
	b.End()
	b.TextElement("button", "切换坐标 / action 锚点", "data-ref", "move", "type", "button")
	return b.String()
}()

const formCSS = `form { display: flex; flex-direction: column; gap: 12px; } input { padding: 8px; } button { padding: 6px 12px; }`

func main() {
	var publisher *peb.ActionPublisher
	name := "Pencode"
	anchored := false
	closeOnOutsideClick := false
	popupPlacement := func() peb.UIPlacement {
		if anchored {
			return peb.UIPlacement{Mode: peb.PositionAnchor, X: -480, Y: 8, Width: 520, Height: 360}
		}
		return peb.UIPlacement{Mode: peb.PositionAbsolute, X: 80, Y: 100, Width: 520, Height: 360}
	}
	panel := func(c peb.ActionContext) error {
		label := name
		if c.Path != "" {
			label += " · " + c.Path
		}
		var b markup.Builder
		b.TextElement("p", label)
		b.TextElement("button", "刷新", "data-ref", "refresh", "type", "button")
		b.TextElement("button", "关闭 Panel", "data-ref", "close", "type", "button")
		return publisher.Show(peb.UI{ID: "result", Title: "UI 示例", Kind: peb.Panel, Context: c, View: &peb.EditorView{
			HTML:   b.String(),
			CSS:    `button { margin-right: 8px; padding: 6px 12px; }`,
			Events: []peb.EditorBinding{{Ref: "refresh", Kind: "click"}, {Ref: "close", Kind: "click"}},
		}})
	}
	err := peb.ServeExtension(os.Stdin, os.Stdout, peb.Contributions{Settings: []peb.Setting{{ID: "name", Label: "默认名称", Description: "在结果 Panel 中显示的默认名称。", Default: json.RawMessage(`"Pencode"`)}}, Actions: []peb.Action{{ID: "tools", Text: "UI 示例", Tooltip: "扩展菜单、弹窗与 Panel", Location: peb.TitleBarRight}}}, peb.Handlers{
		Settings: func(values peb.SettingValues) error {
			var err error
			name, err = peb.DecodeSetting[string](values, "name")
			return err
		},
		ActivateActions: func(_ context.Context, p *peb.ActionPublisher) { publisher = p },
		Action: func(c peb.ActionContext) error {
			return publisher.Show(peb.UI{ID: "menu", Title: "UI 示例操作", Anchor: c.ID, Kind: peb.ContextMenu, View: &peb.EditorView{HTML: menuHTML, Events: []peb.EditorBinding{{Ref: "dialog", Kind: "click"}, {Ref: "popup", Kind: "click"}, {Ref: "panel", Kind: "click"}, {Ref: "notification", Kind: "click"}, {Ref: "auto-close", Kind: "click"}}}})
		},
		UIEvent: func(event peb.UIEvent) error {
			switch event.Event.Kind {
			case "open", "close", "visibilitychange":
				fmt.Fprintf(os.Stderr, "%s: %s %s %s\n", event.ID, event.Event.Kind, event.Reason, event.Event.Value)
			case "input":
				if event.ID == "form" && event.Event.Ref == "name" {
					name = event.Event.Value
				}
				if event.ID == "floating" {
					fmt.Fprintf(os.Stderr, "搜索: %s\n", event.Event.Value)
				}
			case "change":
				if event.ID == "floating" && event.Event.Ref == "outside" {
					closeOnOutsideClick = event.Event.Checked
					return publisher.Show(peb.UI{ID: "floating", Title: "自由定位浮窗", Kind: peb.Popup, Placement: popupPlacement(), CloseOnOutsideClick: closeOnOutsideClick})
				}
			case "submit":
				if event.ID == "form" {
					if err := publisher.CloseUI("form"); err != nil {
						return err
					}
					return panel(event.Context)
				}
			case "click":
				switch event.ID {
				case "menu":
					if event.Event.Ref == "notification" || event.Event.Ref == "auto-close" {
						if err := publisher.CloseUI("menu"); err != nil {
							return err
						}
						if event.Event.Ref == "auto-close" {
							return publisher.Show(peb.UI{ID: "temporary", Title: "操作完成", Kind: peb.Notification, AutoCloseMS: 5000, View: &peb.EditorView{HTML: temporaryHTML}})
						}
						return publisher.Show(peb.UI{ID: "notification", Title: "UI 示例", Kind: peb.Notification, View: &peb.EditorView{
							HTML:   notificationHTML,
							CSS:    `footer { display: flex; justify-content: flex-end; gap: 8px; }`,
							Events: []peb.EditorBinding{{Ref: "result", Kind: "click"}, {Ref: "dismiss", Kind: "click"}},
						}})
					}
					if event.Event.Ref == "popup" {
						anchored = false
						closeOnOutsideClick = false
						return publisher.Show(peb.UI{ID: "floating", Title: "自由定位浮窗", Kind: peb.Popup, Placement: popupPlacement(), View: &peb.EditorView{
							HTML:   popupHTML,
							CSS:    `input[data-ref="query"] { display: block; width: 100%; padding: 8px; margin: 8px 0 16px; } button { padding: 6px 12px; }`,
							Events: []peb.EditorBinding{{Ref: "query", Kind: "input"}, {Ref: "outside", Kind: "change"}, {Ref: "move", Kind: "click"}},
						}})
					}
					if event.Event.Ref == "dialog" {
						return publisher.Show(peb.UI{ID: "form", Title: "输入名称", Kind: peb.Dialog, View: &peb.EditorView{HTML: formHTML, CSS: formCSS, Events: []peb.EditorBinding{{Ref: "name", Kind: "input"}, {Ref: "form", Kind: "submit"}}}})
					}
					if event.Event.Ref == "panel" {
						if err := publisher.CloseUI("menu"); err != nil {
							return err
						}
						return panel(event.Context)
					}
				case "notification":
					if err := publisher.CloseUI(event.ID); err != nil {
						return err
					}
					if event.Event.Ref == "result" {
						return panel(event.Context)
					}
				case "floating":
					if event.Event.Ref == "move" {
						anchored = !anchored
						return publisher.Show(peb.UI{ID: "floating", Title: "自由定位浮窗", Kind: peb.Popup, Placement: popupPlacement(), CloseOnOutsideClick: closeOnOutsideClick})
					}
				case "result":
					if event.Event.Ref == "refresh" {
						return panel(event.Context)
					}
					if event.Event.Ref == "close" {
						return publisher.CloseUI("result")
					}
				}
			}
			return nil
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
