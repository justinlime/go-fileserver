let oList = [];
for (let orig of document.getElementsByClassName('item')) {
   oList.push(orig.style.display) 
}
function search() {
  let input = document.getElementById('searchbar').value
  input = input.toLowerCase();
  let x = document.getElementsByClassName('item');

  for (i = 0; i < x.length; i++) {
    if (!x[i].innerHTML.toLowerCase().includes(input)) {
      x[i].style.display = "none";
    }
    else {
      x[i].style.display = oList[i];
    }
  }
}
