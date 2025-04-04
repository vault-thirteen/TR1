// Settings.
rootPath = "/";
apiPath = "/api";
settingsPath = "/settings";
captchaPath = "/captcha"
redirectDelay = 5;
urlParameter_Id = "id";
urlParameter_Action = "a";

// Action pages.
ActionPage = {
    LogIn: "log-in",
    Register: "register",
}

// HTTP Status Codes.

httpStatusCode_NotAuthorised = 401;

// Messages.
Msg = {
    GenericErrorPrefix: "Error: ",
    Redirecting: "Redirecting. Please wait ...",
}

// Errors.
Err = {
    ActionMismatch: "action mismatch",
    Client: "client error",
    RpcError: "RPC error",
    Server: "server error",
    Settings: "settings error",
    Unknown: "unknown error",
    BooleanToString: "booleanToString: ",
    UnknownActionPage: "Unknown action page: ",
    UnknownElementType: "Unknown element type: ",
    UnknownPageContentType: "Unknown page content type: ",
    EmailAddressIsNotValid: "E-Mail address is not valid",
    CaptchaAnswerIsNotSet: "Captcha answer is not set",
    VerificationCodeIsNotSet: "Verification code is not set",
    PasswordIsNotSet: "Password is not set",
    PasswordIsNotAllowed: "Password is not allowed",
    RequestIdIsNotSet: "Request ID is not set",
    PasswordIsDifferent: "Password is different",
    NameIsNotSet: "Name is not set",
}

// Action names.
ActionName = {
    // AuthService.
    ConfirmLogIn: "confirmLogIn",
    ConfirmRegistration: "confirmRegistration",
    GetSelfRoles: "getSelfRoles",
    StartLogIn: "startLogIn",
    StartRegistration: "startRegistration",

    // MessageService.
    AddForum: "addForum",
    ListForums: "listForums",
}

class ApiRequest {
    constructor(action, parameters) {
        this.action = action;
        this.parameters = parameters;
    }
}

class ApiResponse {
    constructor(isOk, jsonObject, statusCode, errorText) {
        this.isOk = isOk;
        this.jsonObject = jsonObject;
        this.statusCode = statusCode;
        this.errorText = errorText;
    }
}


// Basic API methods.

let lastHttpStatusCode = 0;

async function sendApiRequest(data) {
    let url = apiPath;
    let ri = {
        method: "POST",
        body: JSON.stringify(data)
    };
    let resp = await fetch(url, ri);
    let result;
    lastHttpStatusCode = resp.status;

    if (resp.status === 200) {
        result = new ApiResponse(true, await resp.json(), resp.status, null);
        return result;
    } else {
        result = new ApiResponse(false, null, resp.status, await resp.text());
        if (result.errorText.length === 0) {
            result.errorText = createErrorTextByStatusCode(result.statusCode);
        }
        return result;
    }
}

function createErrorTextByStatusCode(statusCode) {
    if ((statusCode >= 400) && (statusCode <= 499)) {
        return Msg.GenericErrorPrefix + Err.Client + " (" + statusCode.toString() + ")";
    }
    if ((statusCode >= 500) && (statusCode <= 599)) {
        return Msg.GenericErrorPrefix + Err.Server + " (" + statusCode.toString() + ")";
    }
    return Msg.GenericErrorPrefix + Err.Unknown + " (" + statusCode.toString() + ")";
}

async function sendApiRequestAndGetResult(reqData) {
    let actionName = reqData.action;

    let resp = await sendApiRequest(reqData);
    if (!resp.isOk) {
        console.error(composeErrorText(resp.errorText));
        return null;
    }

    let jo = resp.jsonObject;
    if (jo == null) {
        console.error(composeErrorText(Err.RpcError));
        return null;
    }

    if (jo.action !== actionName) {
        console.error(composeErrorText(Err.ActionMismatch));
        return null;
    }

    return jo.result;
}

function composeErrorText(errMsg) {
    return Msg.GenericErrorPrefix + errMsg.trim() + ".";
}

function makeUrl_CaptchaImage(id) {
    return captchaPath + '?' + urlParameter_Id + '=' + id;
}

function makeUrl_ActionPage(action) {
    return rootPath + '?' + urlParameter_Action + '=' + action;
}

// Models.

// N.B. One year ago classes worked fine with object field names starting with
// a capital letter. Somewhere in the past time someone changed this behaviour
// and now object field names are not automatically converted to lower case
// when being used in JSON format. This is a fundamental breaking change in
// JavaScript language. JavaScript must die.

class Forum {
    constructor(id, name, threads, pos) {
        this.id = id;
        this.name = name;
        this.threads = threads;
        this.pos = pos;
    }
}

class User {
    constructor(id, name, email, password, session, roles, regTime, banTime) {
        this.id = id;
        this.name = name;
        this.email = email;
        this.password = new Password(password);
        this.session = session;
        this.roles = roles;
        this.regTime = regTime;
        this.banTime = banTime;
    }
}

class Password {
    constructor(text) {
        this.text = text;
    }
}

// Request parameters & RPC functions.

// AuthService.

class Parameters_ConfirmLogIn {
    constructor(requestId, captchaAnswer, verificationCode, authData) {
        this.requestId = requestId;
        this.captchaAnswer = captchaAnswer;
        this.verificationCode = verificationCode;
        this.authData = authData;
    }
}

async function confirmLogIn(requestId, captchaAnswer, verificationCode, authData) {
    let params = new Parameters_ConfirmLogIn(requestId, captchaAnswer, verificationCode, authData);
    let reqData = new ApiRequest(ActionName.ConfirmLogIn, params);
    return await sendApiRequestAndGetResult(reqData);
}

class Parameters_ConfirmRegistration {
    constructor(requestId, captchaAnswer, verificationCode) {
        this.requestId = requestId;
        this.captchaAnswer = captchaAnswer;
        this.verificationCode = verificationCode;
    }
}

async function confirmRegistration(requestId, captchaAnswer, verificationCode) {
    let params = new Parameters_ConfirmRegistration(requestId, captchaAnswer, verificationCode);
    let reqData = new ApiRequest(ActionName.ConfirmRegistration, params);
    return await sendApiRequestAndGetResult(reqData);
}

async function getSelfRoles() {
    let reqData = new ApiRequest(ActionName.GetSelfRoles, {});
    return await sendApiRequestAndGetResult(reqData);
}

class Parameters_StartLogIn {
    constructor(email) {
        this.user = new User(0, "", email);
    }
}

async function startLogIn(email) {
    let params = new Parameters_StartLogIn(email);
    let reqData = new ApiRequest(ActionName.StartLogIn, params);
    return await sendApiRequestAndGetResult(reqData);
}

class Parameters_StartRegistration {
    constructor(name, email, password) {
        this.user = new User(0, name, email, password);
    }
}

async function startRegistration(name, email, password) {
    let params = new Parameters_StartRegistration(name, email, password);
    let reqData = new ApiRequest(ActionName.StartRegistration, params);
    return await sendApiRequestAndGetResult(reqData);
}

// MessageService.

class Parameters_AddForum {
    constructor(name) {
        this.forum = new Forum(0, name);
    }
}

async function addForum(name) {
    let params = new Parameters_AddForum(name);
    let reqData = new ApiRequest(ActionName.AddForum, params);
    return await sendApiRequestAndGetResult(reqData);
}

async function listForums() {
    let reqData = new ApiRequest(ActionName.ListForums, {});
    return await sendApiRequestAndGetResult(reqData);
}
