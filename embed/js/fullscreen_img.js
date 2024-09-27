// search
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
// images
document.addEventListener("DOMContentLoaded", function() {
  // Select all image elements
  const images = document.querySelectorAll('img');
    
  // Define the function to be executed on click
  function handleImageClick(event) {
    const img = event.target;
    // Open the image in fullscreen
    openFullscreen(img);
  }
    
  // Add the click and touch event listener to each image
  images.forEach(img => {
    img.addEventListener('click', handleImageClick);
    // img.addEventListener('touchend', handleImageClick); // Handle touch events
  });
    
  // Function to handle fullscreen view
  function openFullscreen(img) {
    // Create and display a modal for fullscreen view
    const modal = document.createElement('div');
    modal.style.position = 'fixed';
    modal.style.zIndex = '1000';
    modal.style.left = '0';
    modal.style.top = '0';
    modal.style.width = '100%';
    modal.style.height = '100%';
    
    modal.style.background = "rgba(0,0,0,0.6)"
    modal.style.backdropFilter = "blur(8px)"
    modal.style.display = 'flex';
    modal.style.justifyContent = 'center';
    modal.style.alignItems = 'center';
    modal.style.cursor = 'pointer'; // Make sure clicks/taps are captured

    const fullscreenImg = document.createElement('img');
    const widthScale = window.innerWidth / img.naturalWidth; 
    const heightScale = window.innerHeight / img.naturalHeight; 
    let scaleFactor = Math.min(widthScale, heightScale) * 0.95;
    if (scaleFactor < 1) {
      scaleFactor = 1
    }
    fullscreenImg.src = img.src;
    fullscreenImg.style.transform = `scale(${scaleFactor})`;
    fullscreenImg.style.maxWidth = '95%';
    fullscreenImg.style.maxHeight = '95%';
    fullscreenImg.style.margin = '0';
    fullscreenImg.style.display = 'block'; // Ensure image is displayed correctly
    fullscreenImg.style.borderRadius = '0';

    const closeBtn = document.createElement('span');
    closeBtn.textContent = '×';
    closeBtn.style.position = 'absolute';
    closeBtn.style.top = '15px';
    closeBtn.style.right = '15px';
    closeBtn.style.color = '#F38BA8';
    closeBtn.style.fontSize = '40px';
    closeBtn.style.fontWeight = 'bold';
    closeBtn.style.cursor = 'pointer';
        
    modal.addEventListener('click', () => {
      document.body.removeChild(modal);
    });

    document.addEventListener("keydown", event => {
      if (event.key === "Escape" && document.contains(modal)) {
        document.body.removeChild(modal);
      }
    })

    // Append image and close button to the modal
    modal.appendChild(fullscreenImg);
    modal.appendChild(closeBtn);
    document.body.appendChild(modal);
  }
});


