
console.log("Script linked successfully.");
console.log(mdc);

const topAppBarElement = document.querySelector('.mdc-top-app-bar');
const topAppBar = mdc.topAppBar.MDCTopAppBar.attachTo(topAppBarElement);
const drawer = mdc.drawer.MDCDrawer.attachTo(document.querySelector('.mdc-drawer'));

topAppBar.setScrollTarget(document.getElementById('scrollbar'));
topAppBar.listen('MDCTopAppBar:nav', () => {
    drawer.open = !drawer.open;
});

if ($('#sessionUser').text() != "") {
    const menu = mdc.menu.MDCMenu.attachTo(document.querySelector('.mdc-menu'));
    const dashButton = document.querySelector('#dashboard');

    dashButton.addEventListener('click', () => {
        menu.open = !menu.open;
    });
}