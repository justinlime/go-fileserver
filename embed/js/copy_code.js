const highlightDiv = document.querySelector(".text-preview-container")

const button = document.createElement("button");
const preElement = highlightDiv.querySelector('pre');
const code = preElement.querySelector('code');
highlightDiv.style.position = "relative";
button.className = "copy-code-button";
button.type = "button";
button.innerText = "Copy";
button.style.opacity = "0";
button.style.transition = "opacity 0.25s ease-in-out";
button.style.position = "absolute";
button.style.cursor = 'pointer';
button.addEventListener("click", () => {
  const range = document.createRange();
  range.selectNode(code);
  window.getSelection().removeAllRanges();
  window.getSelection().addRange(range);
  if (navigator.clipboard) {
    navigator.clipboard.writeText(window.getSelection().toString());
  } else {
    document.execCommand('copy')
  }
  window.getSelection().removeAllRanges();
  button.innerHTML = "Copied!";
  setTimeout(() => {
    button.innerHTML = "Copy";
  }, 2000)
})
preElement.appendChild(button);

const butt = highlightDiv.querySelector('button');
highlightDiv.addEventListener('mouseover', () => {
  // butt.style.display = 'inline-block';
  butt.style.opacity = "1";
}); 
highlightDiv.addEventListener('mouseout', () => {
  // butt.style.display = 'none';
  butt.style.opacity = "0";
});
// for mobile 
window.addEventListener('scroll', () => {
  // butt.style.display = 'none'; 
  butt.style.opacity = "0";
});
