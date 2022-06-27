package core

const FireEvents = `
function fireEvent(element,event){
    var evt = document.createEvent("HTMLEvents");
    evt.initEvent(event, true, true );
    return !element.dispatchEvent(evt);
}
var els = document.getElementsByTagName("select");
for(var i=0; i<els.length; i++) {
    fireEvent(els[i],'change');
    fireEvent(els[i],'click');
}
var els = document.getElementsByTagName("input");
for(var i=0; i<els.length; i++) {
    fireEvent(els[i],'change');
    fireEvent(els[i],'click');
}`

const FillAndSubmitForms = `
getForms();
function getForms(){
    var formarray = [];
    var forms = document.getElementsByTagName("form");
    for(var i=0; i<forms.length; i++) {
        try {
            var meth = "";
            if ( forms[i].method.toUpperCase() != "GET" ){
                meth = "#"+forms[i].method.toUpperCase()
            }
            forms[i].target = '_blank';
            var inputsarray = [];
            var inputs = forms[i].getElementsByTagName("input");
            for(var a=0; a<forms.length; a++) {
                if (inputs[a].type != 'submit'){
                    inputs[a].value = 'FUZZ';
                }
            }
            forms[i].submit();
        }catch(e){}
    }
    return formarray
}`

const GetForms = `
getForms();
function getForms(){
    var formarray = [];
    var forms = document.getElementsByTagName("form");
    for(var i=0; i<forms.length; i++) {
        try {var meth = "";
        if ( forms[i].method.toUpperCase() != "GET" ){
            meth = "&is"+forms[i].method.toUpperCase()
        }
        var inputsarray = [];
        var inputs = forms[i].getElementsByTagName("input");
        for(var a=0; a<forms.length; a++) {
            if (inputs[a].type != 'submit'){
                inputsarray.push(inputs[a].name+"="+inputs[a].value)
            }
        }
        formarray.push(forms[i].action + "?"+inputsarray.join('&')+meth);
    }catch(e){}
}
return formarray
}`

const GetLinks = `
getLinks();
function absolutePath(href) {
    try {
        var link = document.createElement("a");
        link.href = href;
        return link.href;
    } catch (error) {}
}
function getLinks() {
    var array = [];
    if (!document) return array;
    var allElements = document.querySelectorAll("*");
    for (var el of allElements) {
        if (el.href && typeof el.href === 'string') {
            array.push(el.href);
        } else if (el.src && typeof el.src === 'string') {
            var absolute = absolutePath(el.src);
            array.push(absolute);
        }
    }
    return array;
}`
