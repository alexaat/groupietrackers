function sort(event){

//block 8, 37, 38, 39, 40   
if(event.which == 8 || event.which == 37 || event.which == 38 || event.which == 39 || event.which == 40){
   return
}

let box = document.getElementById("search");
let searchValue = box.value;
let dataList = document.getElementById("artists");
var optionsTop = '';
var optionsBottom = '';
   for (let i = 0; i<dataList.options.length; i++){
      let optionValue = dataList.options[i].value;
      let optionValueLower = optionValue.toLowerCase();
      let result = optionValueLower.startsWith(searchValue.toLowerCase());
      if(result){
            console.log(optionValue);
            optionsTop += '<option value="' + optionValue + '" />';
      }else{
         optionsBottom+='<option value="' + optionValue + '" />';
      }   
   }
   let options = optionsTop + optionsBottom;
   dataList.innerHTML = options;
}

