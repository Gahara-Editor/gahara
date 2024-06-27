export function scrollVertical(node: HTMLLIElement) {
  const listContainer = document.getElementById("content-wrap");
  const listRect = listContainer.getBoundingClientRect();
  const nodeRect = node.getBoundingClientRect();

  const isNodeVisible =
    nodeRect.top >= listRect.top && nodeRect.bottom <= listRect.bottom;

  if (!isNodeVisible) {
    const scrollY = nodeRect.top - listRect.top + listContainer.scrollTop;
    listContainer.scrollTo({
      top: scrollY,
      behavior: "smooth",
    });
  }
}
