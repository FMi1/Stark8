// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Home() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = headerComponent().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<body class=\"min-h-screen bg-gradient-to-b\"><div class=\"p-1 flex flex-wrap items-center justify-center\"><a class=\"cursor-pointer flex-shrink-0 m-4 relative overflow-hidden rounded-lg shadow-sky-300/50 shadow-2xl group\" style=\"width: 144px; height: 155px;\" onclick=\"my_modal_1.showModal()\" hx-get=\"/new\" hx-trigger=\"click\" hx-target=\"#my_modal_1\" hx-swap=\"innerHTML\"><div class=\"flex h-full flex-row justify-center items-center z-30\"><svg xmlns=\"http://www.w3.org/2000/svg\" width=\"36px\" height=\"36px\" viewBox=\"0 0 24 24\" class=\"stroke-zinc-400 fill-none group-hover:fill-zinc-800 group-active:stroke-zinc-200 group-active:fill-zinc-600 group-active:duration-0 duration-300\"><path d=\"M12 22C17.5 22 22 17.5 22 12C22 6.5 17.5 2 12 2C6.5 2 2 6.5 2 12C2 17.5 6.5 22 12 22Z\" stroke-width=\"1.5\"></path> <path d=\"M8 12H16\" stroke-width=\"1.5\"></path> <path d=\"M12 16V8\" stroke-width=\"1.5\"></path></svg> <span class=\"absolute w-full h-full top-0 left-0 bg-sky-50 rounded-md transform scale-x-0 group-hover:scale-x-100 transition-transform group-hover:duration-500 duration-1000 origin-center\"></span> <span class=\"absolute w-full h-full top-0 left-0 bg-sky-100 rounded-md transform scale-x-0 group-hover:scale-x-100 transition-transform group-hover:duration-700 duration-700 origin-center\"></span> <span class=\"absolute w-full h-full top-0 left-0 bg-sky-200 rounded-md transform scale-x-0 group-hover:scale-x-100 transition-transform group-hover:duration-1000 duration-500 origin-center\"></span> <span class=\"group-hover:opacity-100 group-hover:duration-1000 duration-100 opacity-0 absolute z-10 font-semibold text-md text-neutral-500\">Create a Stark8!</span></div></a> <dialog id=\"my_modal_1\" class=\"modal modal-top py-16 flex justify-center rounded-md\"></dialog><div class=\"contents\"><div id=\"stark8\" class=\"flex flex-wrap items-center justify-center\" hx-get=\"/stark8s\" hx-trigger=\"load\"><!-- Your Stark8 items here --></div></div></div><div id=\"toast\" class=\"toast toast-end hidden\"><div class=\"alert alert-success\"><span>Stark8 create successfully.</span></div></div></body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
