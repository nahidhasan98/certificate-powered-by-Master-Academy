// mdc ripple effect
const snackbar = mdc.snackbar.MDCSnackbar.attachTo(document.querySelector('.mdc-snackbar'));

const textFieldElements = [].slice.call(document.querySelectorAll('.mdc-text-field'));
textFieldElements.forEach((textFieldEl) => {
    mdc.textField.MDCTextField.attachTo(textFieldEl);
});

$('#applyForm').submit(function (e) {
    e.preventDefault();

    // checking if any field is empty(only space) or not
    let flag = 1;
    if ($('#studentsName').val().trim() == "") {
        $('#studentsName').val("");
        notify("Please give your full Name!", 5000);
        flag = 0;
    }
    if ($('#codewarsUsername').val().trim() == "") {
        $('#codewarsUsername').val("");
        notify("Please give your codewars username!", 5000);
        flag = 0;
    }

    if (flag) {
        console.log($('#applyForm').serialize());
        $('#btnApply').prop('disabled', true);
        $('#btnApply .mdc-button__label').text("Please wait...");

        //sending ajax post request
        let request = $.ajax({
            async: true,
            type: "POST",
            url: "/apply",
            data: $('#applyForm').serialize(),
        });

        request.done(function (response) {
            notify(response, 5000);

            setTimeout(function () {
                location.reload();
            }, 3000);
        });

        request.fail(function (response) {
            notify(response, 5000);
        });

        request.always(function () {
            $('#btnApply').prop('disabled', false);
            $('#btnApply .mdc-button__label').text("Apply");
        });
    }
    return false;
});

function notify(msg, time) {
    snackbar.timeoutMs = time;
    snackbar.labelText = msg;
    snackbar.actionButtonText = "OKAY";
    snackbar.open();
}