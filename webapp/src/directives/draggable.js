let isDragging = false;
let mouseDown = false;
let currentlyDragging;
let lastMovePos = { x: 0, y: 0 };

let onMouseDown = function(event, el) {
  lastMovePos.x = event.clientX;
  lastMovePos.y = event.clientY;
  currentlyDragging = el;
  mouseDown = true;
};

let onMouseUp = function() {
  isDragging = false;
  mouseDown = false;
};

let onMouseMove = function(event) {
  if (!mouseDown) return;

  if (!isDragging) {
    isDragging = true;
    return;
  }

  let container = document.getElementById(currentlyDragging.dataset.ref);

  let posX = lastMovePos.x - event.clientX;
  let posY = lastMovePos.y - event.clientY;
  lastMovePos.x = event.clientX;
  lastMovePos.y = event.clientY;

  let newX = container.offsetLeft - posX;
  let newY = container.offsetTop - posY;

  if (newX < 0) {
    newX = 0;
  }

  if (newX + container.offsetWidth > document.body.clientWidth) {
    newX = document.body.clientWidth - container.offsetWidth;
  }

  if (newY < 0) {
    newY = 0;
  }

  if (newY + container.offsetHeight > document.body.clientHeight) {
    newY = document.body.clientHeight - container.offsetHeight;
  }

  container.style.left = newX + "px";
  container.style.top = newY + "px";
};

document.addEventListener("mousemove", onMouseMove.bind(this));
document.addEventListener("mouseup", onMouseUp);

export default {
  mounted: function(el) {
    el.addEventListener("mousedown", event => {
      onMouseDown(event, el);
    });
  }
};
