function toggle_advsearch_pane(b) {
    b.innerHTML == "Hide Advanced Options" ? hide_advsearch_pane(b) : show_advsearch_pane(b)
}
function show_advsearch_pane(b) {
    var c = document.getElementById("advdiv");
    b.innerHTML = "Hide Advanced Options";
    c.style.display = "";
    c.innerHTML = '<input type="hidden" id="advsearch" name="advsearch" value="1" /><table class="itss">	<tr>		<td class="ic4"><input id="adv11" type="checkbox" name="f_sname" checked="checked" /> <label for="adv11">Search Gallery Name</label></td>		<td class="ic4"><input id="adv12" type="checkbox" name="f_stags" checked="checked" /> <label for="adv12">Search Gallery Tags</label></td>		<td class="ic2"><input id="adv13" type="checkbox" name="f_sdesc" colspan="2" /> <label for="adv13">Search Gallery Description</label></td>	</tr>	<tr>		<td class="ic2" colspan="2"><input id="adv15" type="checkbox" name="f_storr" /> <label for="adv15">Search Torrent Filenames</label></td>		<td class="ic2" colspan="2"><input id="adv16" type="checkbox" name="f_sto" /> <label for="adv16">Only Show Galleries With Torrents</label></td>	</tr>	<tr>		<td class="ic2" colspan="2"><input id="adv21" type="checkbox" name="f_sdt1" /> <label for="adv21">Search Low-Power Tags</label></td>		<td class="ic2" colspan="2"><input id="adv22" type="checkbox" name="f_sdt2" /> <label for="adv22">Search Downvoted Tags</label></td>	</tr>	<tr>		<td class="ic2" colspan="2"><input id="adv31" type="checkbox" name="f_sh" /> <label for="adv31">Show Expunged Galleries</label></td>		<td class="ic2" colspan="2"><input id="adv32" type="checkbox" name="f_sr" /> <label for="adv32">Minimum Rating:</label> <select id="adv42" class="imr" name="f_srdd"><option value="2">2 stars</option><option value="3">3 stars</option><option value="4">4 stars</option><option value="5">5 stars</option></select></td>	</tr>	<tr>		<td class="ic2" colspan="2"><input id="adv61" type="checkbox" name="f_sp" /> <label for="adv61">Between </label> <input type="text" name="f_spf" value="" size="4" maxlength="4" style="width:30px" /> and <input type="text" name="f_spt" value="" size="4" maxlength="4" style="width:30px" /> pages</td>	</tr>	<tr>		<td class="ic1" colspan="4">Disable default filters for: <input id="adv51" type="checkbox" name="f_sfl" /> <label for="adv51">Language</label> <input id="adv52" type="checkbox" name="f_sfu" /> <label for="adv52">Uploader</label> <input id="adv53" type="checkbox" name="f_sft" /> <label for="adv53">Tags</label></td>	</tr></table>'
}
function hide_advsearch_pane(b) {
    var c = document.getElementById("advdiv");
    b.innerHTML = "Show Advanced Options";
    c.style.display = "none";
    c.innerHTML = ""
}
function load_pane_image(c) {
    if (c != undefined) {
        var a = c.childNodes[0].childNodes[0];
        var b = a.getAttribute("data-src");
        if (b != undefined) {
            a.src = b;
            a.removeAttribute("data-src")
        }
    }
}
function preload_pane_image(b, a) {
    setTimeout(function() {
        if (b > 0) {
            load_pane_image(document.getElementById("it" + b))
        }
        if (a > 0) {
            load_pane_image(document.getElementById("it" + a))
        }
    }, 100)
}
var visible_pane = 0;
function show_image_pane(a) {
    if (visible_pane > 0) {
        hide_image_pane(visible_pane)
    }
    var b = document.getElementById("it" + a);
    load_pane_image(b);
    b.style.visibility = "visible";
    document.getElementById("ic" + a).style.visibility = "visible";
    visible_pane = a
}
function hide_image_pane(a) {
    document.getElementById("it" + a).style.visibility = "hidden";
    document.getElementById("ic" + a).style.visibility = "hidden";
    visible_pane = 0
}
function update_favsel(b) {
    if ((b.value).match(/^fav([0-9])$/)) {
        var a = parseInt(b.value.replace("fav", ""));
        b.style.paddingLeft = "20px";
        b.style.backgroundPosition = "4px " + -(-2 + a * 19) + "px"
    } else {
        b.style.paddingLeft = "2px";
        b.style.backgroundPosition = "4px 20px"
    }
}
function toggle_category(b) {
    var a = document.getElementById("f_cats");
    var c = document.getElementById("cat_" + b);
    if (a.getAttribute("disabled")) {
        a.removeAttribute("disabled")
    }
    if (c.getAttribute("data-disabled")) {
        c.removeAttribute("data-disabled");
        a.value = parseInt(a.value) & (1023 ^ b)
    } else {
        c.setAttribute("data-disabled", 1);
        a.value = parseInt(a.value) | b
    }
}
function search_presubmit() {
    var a = document.getElementById("f_search");
    if (!a.value) {
        a.setAttribute("disabled", 1)
    }
    var c = document.getElementById("adv32");
    var b = document.getElementById("adv42");
    if (!c.checked) {
        b.setAttribute("disabled", 1)
    }
}
function cancel_event(a) {
    a = a ? a : window.event;
    if (a.stopPropagation) {
        a.stopPropagation()
    }
    if (a.preventDefault) {
        a.preventDefault()
    }
    a.cancelBubble = true;
    a.cancel = true;
    a.returnValue = false;
    return false
}
;