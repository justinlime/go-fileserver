function addCopyButton() {
  const highlightDiv = document.querySelector(".text-preview-container")
  if (highlightDiv === null ) {
    return 
  }
  const button = document.createElement("button");
  const preElement = highlightDiv.querySelector('pre');
  const code = preElement.querySelector('code');
  highlightDiv.style.position = "relative";
  button.className = "copy-code-button";
  button.type = "button";
  button.innerText = "Copy";
  button.style.opacity = "0";
  button.style.top = "30px";
  button.style.right = "35px";
  button.style.transition = "opacity 0.25s ease-in-out";
  button.style.position = "absolute";
  button.style.cursor = 'pointer';
  button.style.border = 'none';
  button.style.fontSize = '1.05rem';
  button.style.fontFamily = 'firasans';
  button.style.borderRadius = '0.25rem';
  button.style.backgroundColor = '#585B70';
  button.style.color = '#cdd6f4';
  button.style.padding = '0.5rem 1rem';

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
}
