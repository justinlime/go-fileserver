function addCopy() {
  const container = document.querySelector(".text-preview-container")
  const button = container.querySelector(".copy-button")
  const code = container.querySelector('code');

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

  container.addEventListener('mouseover', () => {
    // butt.style.display = 'inline-block';
    button.style.opacity = "1";
  }); 
  container.addEventListener('mouseout', () => {
    // butt.style.display = 'none';
    button.style.opacity = "0";
  });
  // for mobile 
  window.addEventListener('scroll', () => {
    // butt.style.display = 'none'; 
    button.style.opacity = "0";
  });
}

addCopy()
