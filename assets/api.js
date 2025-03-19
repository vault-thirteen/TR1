// Settings.
rootPath = "/";
apiPath = "/api";
settingsPath = "/settings";
redirectDelay = 3;

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
}

// Action names.
ActionName = {
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

async function sendApiRequest(data) {
    let url = apiPath;
    let ri = {
        method: "POST",
        body: JSON.stringify(data)
    };
    let resp = await fetch(url, ri);
    let result;
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

// Request parameters & RPC functions.

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
