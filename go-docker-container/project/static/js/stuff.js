function buttonAdder(){
    //JSON for button adding
    var ButtonAdd = {
        AddAmount: 1
    };
    var jsonString = JSON.stringify(ButtonAdd); //Stringify Data
    //Send Request to end game
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/addClick', true);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.addEventListener('readystatechange', function(){
        if(xhr.readyState === XMLHttpRequest.DONE && xhr.status === 200){
            var item = xhr.responseText;
            var ReturnData = JSON.parse(item);
            if (ReturnData.SuccessInt == 0){
                console.log("DEBUG: Successful button click");
                //window.location.assign("/gamepage");
            } else {
                console.log("Unsuccessful button click.");
            }
        }
    });
    xhr.send(jsonString);
}