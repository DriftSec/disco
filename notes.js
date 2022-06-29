

// get elements by tagname with javascript
var array = [];
var links = document.getElementsByTagName("a");
for(var i=0; i<links.length; i++) {
    array.push(links[i].href);
}
links[0].href



// list all events (onclick etc) may not work with jquery types
// https://gist.github.com/dmnsgn/36b26dfcd7695d02de77f5342b0979c7
// https://gist.github.com/tkafka/1c5174ed5b446de7dfa4c04a3b09c95f
// https://www.perimeterx.com/tech-blog/2019/list-every-event-that-exists-in-the-browser/ << works but too much
  
//   function _getEvents(obj) {
//     var result = [];
  
//     for (var prop in obj) {
//       if (0 == prop.indexOf("on")) {
//         prop = prop.substr(2); // remove "on" at the beginning
//         result.push(prop);
//       }
//     }
  
//     return result;
//   }
//   function getEvents() {
//     const result = {};
  
//     result["window"] = _getEvents(window, hasOwnProperty);
  
//     const arr = Object.getOwnPropertyNames(window);
  
//     for (let i = 0; i < arr.length; i++) {
//       const element = arr[i];
  
//       let resultArray = [];
  
//       try {
//         const obj = window[element];
  
//         if (!obj || !obj["prototype"]) {
//           continue;
//         }
  
//         proto = obj["prototype"];
  
//         resultArray = _getEvents(proto);
//       } catch (err) {
//         // console.error(`failed to get events of %o`, element);
//       }
  
//       result[element] = resultArray;
//     }
  
//     return result;
//   }


  //get all forms << works
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
        var inputsarray = [];
        var inputs = forms[i].getElementsByTagName("input");
        for(var a=0; a<forms.length; a++) {
            if (inputs[a].type != 'submit'){
                inputsarray.push(inputs[a].name+"="+inputs[a].value)
            }
        }

        formarray.push(forms[i].action + "?"+inputsarray.join('&')+meth);
    }
    catch(e){}
}
    return formarray
}


// fill and submit
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
    //   formarray.push(forms[i].action + "?"+inputsarray.join('&')+meth);
    }
    catch(e){}
}
    return formarray
}

// // request all forms
// getForms();
// async function postData(url = '', data = '') {
//     const response = await fetch(url, {
//       method: 'POST',
//       body: data
//     });
//     return response.text();
//   };
// function getForms(){
//     var formarray = [];
//     var forms = document.getElementsByTagName("form");
//     for(var i=0; i<forms.length; i++) {
//        try {
//         var meth = "";
//         if ( forms[i].method.toUpperCase() != "GET" ){
//             meth = "#"+forms[i].method.toUpperCase()
//         }
//         var inputsarray = [];
//         var inputs = forms[i].getElementsByTagName("input");
//         for(var a=0; a<forms.length; a++) {
//             inputsarray.push(inputs[a].name+"="+inputs[a].value)
//         }
        
//         postData(forms[i].action,inputsarray.join('&'))

//         // formarray.push(forms[i].action + "?"+inputsarray.join('&')+meth);
//     }
//     catch(e){}
// }
//     // return formarray
// }


// fire event works ***************************** goods here
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

}
