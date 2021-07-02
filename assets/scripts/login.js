// mdc ripple effect
mdc.textField.MDCTextField.attachTo(document.querySelector('.mdc-text-field-username'));
mdc.textField.MDCTextField.attachTo(document.querySelector('.mdc-text-field-password'));
mdc.ripple.MDCRipple.attachTo(document.querySelector('.cancel'));
mdc.ripple.MDCRipple.attachTo(document.querySelector('.next'));
const snackbar = mdc.snackbar.MDCSnackbar.attachTo(document.querySelector('.mdc-snackbar'));

$('#loginForm').submit(function (e) {
    e.preventDefault();

    // checking if any field is empty(only space) or not
    let flag = 1;
    if ($('#username').val().trim() == "") {
        $('#username').val("");
        notify("Please give a username!", 5000);
        flag = 0;
    }

    if (flag) {
        $('#btnLogin').prop('disabled', true);
        $('#btnLogin .mdc-button__label').text("Please wait...");

        //sending ajax post request
        let request = $.ajax({
            async: true,
            type: "POST",
            url: "/login",
            data: $('#loginForm').serialize(),
        });

        request.done(function (response) {
            //console.log(response.trim())
            if (response.trim() == "Login ADMIN") {
                notify("Login successful!", 5000);
                window.location.href = "/dashboard";
            } else if (response.trim() == "Login Done") {
                notify("Login successful!", 5000);
                window.location.href = "/";
            } else {
                //console.log(response)
                notify(response, 5000);
            }
        });

        request.fail(function (response) {
            notify(response, 5000);
        });

        request.always(function () {
            $('#btnLogin').prop('disabled', false);
            $('#btnLogin .mdc-button__label').text("Login");
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