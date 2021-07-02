//js instantiate for mdc elements
const dialog = mdc.dialog.MDCDialog.attachTo(document.querySelector('.mdc-dialog'));
const dialogSearch = mdc.dialog.MDCDialog.attachTo(document.querySelector('#searchFormDialog'));
const snackbar = mdc.snackbar.MDCSnackbar.attachTo(document.querySelector('.mdc-snackbar'));

const textFieldElements = [].slice.call(document.querySelectorAll('.mdc-text-field'));
textFieldElements.forEach((textFieldEl) => {
    mdc.textField.MDCTextField.attachTo(textFieldEl);
});

$('#btnFilter').on('click', function (e) {
    console.log("Button filter pressed.");
    //notify("test", 6000)

    dialogSearch.autoStackButtons = false;
    dialogSearch.open();
    dialogSearch.escapeKeyAction = "";
    dialogSearch.scrimClickAction = "";

    $('#searchForm').trigger("reset");
});

$('#btnSend').on('click', function () {
    console.log("btnSend pressed.");
    notify("Sending mail. Please wait for the confirmation message.", 10000);

    let arr = [];
    for (let item of emailList) {
        arr.push(item)
    }

    //retrieving data from backend
    //sending ajax post request
    let request = $.ajax({
        async: true,
        type: "POST",
        url: "/api/email",
        data: { list: JSON.stringify(arr) },
    });
    request.done(function (response) {
        //console.log(response["Success"]);
        notify("Email sent successfully.", 5000)

        setTimeout(function () {
            location.reload();
        }, 2000);
    });

    request.fail(function (response) {
        notify(response, 5000);
    });
});

function notify(msg, time) {
    snackbar.timeoutMs = time;
    snackbar.labelText = msg;
    snackbar.actionButtonText = "OKAY";
    snackbar.open();
}

function certify(cr) {
    notify("Processing. Please wait...", 10000)
    selectedRow = cr.parentElement.parentElement;
    let email = selectedRow.cells[2].innerText.trim();

    //retrieving data from backend
    //sending ajax post request
    let request = $.ajax({
        async: true,
        type: "GET",
        url: "/api/pdf-" + email,
    });
    request.done(function (response) {
        //console.log(response);
        window.open(response, '_blank');
    });

    request.fail(function (response) {
        notify(response, 5000);
    });
}

//checkbox clicking
let emailList = new Set();
$('.mdc-data-table__table > thead > tr').find('input[type=checkbox]').on('change', function () {
    let isChecked = $(this).is(':checked');


    $('.mdc-data-table__table > tbody  > tr').each(function (index, tr) {
        $(tr).find('input[type=checkbox]').eq(0).prop('checked', isChecked);

        let email = tr.cells[2].innerText.trim();

        if (isChecked == true) {
            emailList.add(email);
        } else {
            emailList.clear();
        }
    });

    //show/hide send button
    if (emailList.size > 0) {
        $('#btnSend').show();
    } else {
        $('#btnSend').hide();
    }

    //console.log(emailList);
});

$('.mdc-data-table__table > tbody  > tr').on("click", ".mdc-checkbox__native-control", function () {
    let isChecked = $(this).is(':checked');
    selectedRow = $(this).closest("td").closest('tr');
    let email = selectedRow[0].cells[2].innerText.trim();

    if (isChecked == true) {
        emailList.add(email);
    } else {
        emailList.delete(email);
    }

    //show/hide send button
    if (emailList.size > 0) {
        $('#btnSend').show();
    } else {
        $('#btnSend').hide();
    }

    //console.log(emailList);
});