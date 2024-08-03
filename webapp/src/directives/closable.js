export default {
  mounted: function(el, binding, vnode) {
    el.clickOutsideEvent = function(event) {
      if (!(el == event.target || el.contains(event.target))) {
        vnode.context[binding.expression](event);
      }
    };
    document.body.addEventListener("click", el.clickOutsideEvent);
  },
  unmounted: function(el) {
    document.body.removeEventListener("click", el.clickOutsideEvent);
  }
}
