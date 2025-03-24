window.onpageshow = function (event) {
    if (event.persisted) {
        // Unfortunately, JavaScript does not reload a page when you click
        // "Go Back" button in your web browser. Every year programmers invent
        // a new "wheel" to fix this bug. And every year old working solutions
        // stop working and new ones are invented. This circus looks infinite,
        // but in reality it will end as soon as this evil programming language
        // dies. Please, do not support JavaScript and its developers in any
        // means possible. Please, let this evil "technology" to die.
        console.info("JavaScript must die. This pseudo language is a big mockery and ridicule of people. This is not a joke. This is truth.");
        window.location.reload();
    }
};

// Names of JavaScript storage variables.
Varname = {
    // Settings.
    Settings_LoadTime: "settings_LoadTime",
    Settings_Version: "settings_Version",
    Settings_TTL: "settings_TTL",
    Settings_SiteName: "settings_SiteName",
    Settings_SiteDomain: "settings_SiteDomain",
    Settings_SessionMaxDuration: "settings_SessionMaxDuration",
    Settings_MessageEditTime: "settings_MessageEditTime",
    Settings_PageSize: "settings_PageSize",
}

// Page sections.
id_td_pageHeader = "pageHeader";
id_td_pageContent = "pageContent";
id_td_pageFooter = "pageFooter";

// IDs of various sections.
id_table_logIn = "logIn";

// Page content types.
pageContent_logIn = "logIn"

// Common basic functions.

function isNumber(x) {
    return typeof x === 'number';
}

function isNumeric(str) {
    if (typeof str != "string") {
        return false
    }

    return !isNaN(str) && !isNaN(parseFloat(str))
}

function booleanToString(b) {
    if (b === true) {
        return "Yes";
    }
    if (b === false) {
        return "No";
    }
    console.error("booleanToString:", b);
    return null;
}

function stringToBoolean(s) {
    if (s == null) {
        return null;
    }

    let x = s.trim().toLowerCase();

    switch (x) {
        case "true":
            return true;

        case "false":
            return false;

        case "yes":
        case "1":
            return true;

        case "no":
        case "0":
            return false;

        default:
            return JSON.parse(x);
    }
}

function getCurrentTimestamp() {
    return Math.floor(Date.now() / 1000);
}

async function sleep(ms) {
    await new Promise(r => setTimeout(r, ms));
}

function addTimeSec(t, deltaSec) {
    return new Date(t.getTime() + deltaSec * 1000);
}

function prettyTime(timeStr) {
    if (timeStr == null) {
        return "";
    }
    if (timeStr.length === 0) {
        return "";
    }

    let t = new Date(timeStr);
    let monthN = t.getUTCMonth() + 1; // Months in JavaScript start with 0 !

    return t.getUTCDate().toString().padStart(2, '0') + "." +
        monthN.toString().padStart(2, '0') + "." +
        t.getUTCFullYear().toString().padStart(4, '0') + " " +
        t.getUTCHours().toString().padStart(2, '0') + ":" +
        t.getUTCMinutes().toString().padStart(2, '0');
}

function countSymbol(str, symbol) {
    let x = str.replaceAll(symbol, '');
    return str.length - x.length;
}

function validateEmailAddress(x) {
    if (typeof x !== 'string') {
        return false
    }
    if (x.length < 3) {
        return false
    }
    if (countSymbol(x, '@') !== 1) {
        return false;
    }
    let atPos = x.indexOf('@');
    if ((atPos === 0) || (atPos === x.length - 1)) {
        return false
    }
    return true;
}

async function redirectPage(wait, url) {
    if (wait) {
        await sleep(redirectDelay * 1000);
    }

    document.location.href = url;
}


// Settings.

class Settings {
    constructor(version, ttl, siteName, siteDomain, sessionMaxDuration, messageEditTime, pageSize) {
        this.Version = version;
        this.TTL = ttl;
        this.SiteName = siteName;
        this.SiteDomain = siteDomain;
        this.SessionMaxDuration = sessionMaxDuration;
        this.MessageEditTime = messageEditTime;
        this.PageSize = pageSize;
    }
}

async function updateSettingsIfNeeded() {
    if (isSettingsUpdateNeeded()) {
        return await updateSettings();
    }
    return true;
}

function isSettingsUpdateNeeded() {
    let settingsLoadTimeStr = sessionStorage.getItem(Varname.Settings_LoadTime);
    if (settingsLoadTimeStr == null) {
        return true;
    }

    let settingsTtlStr = sessionStorage.getItem(Varname.Settings_TTL);
    if (settingsTtlStr == null) {
        return true;
    }
    let settingsTtl = Number(settingsTtlStr);

    let timeNow = getCurrentTimestamp();
    let settingsAge = timeNow - Number(settingsLoadTimeStr);
    if (settingsAge >= settingsTtl) {
        return true;
    }

    return false;
}

async function updateSettings() {
    let resp = await fetchSettings();
    let s = jsonToSettings(resp);
    console.info('New settings have been received. Version: ' + s.Version.toString() + ".");

    // Save the settings for future usage.
    saveSettings(s);
    return true;
}

async function fetchSettings() {
    let data = await fetch(settingsPath);
    return await data.json();
}

function jsonToSettings(x) {
    return new Settings(
        x.version,
        x.ttl,
        x.siteName,
        x.siteDomain,
        x.sessionMaxDuration,
        x.messageEditTime,
        x.pageSize,
    );
}

function saveSettings(s) {
    sessionStorage.setItem(Varname.Settings_Version, s.Version);
    sessionStorage.setItem(Varname.Settings_TTL, s.TTL);
    sessionStorage.setItem(Varname.Settings_SiteName, s.SiteName);
    sessionStorage.setItem(Varname.Settings_SiteDomain, s.SiteDomain);
    sessionStorage.setItem(Varname.Settings_SessionMaxDuration, s.SessionMaxDuration.toString());
    sessionStorage.setItem(Varname.Settings_MessageEditTime, s.MessageEditTime.toString());
    sessionStorage.setItem(Varname.Settings_PageSize, s.PageSize.toString());

    let timeNow = getCurrentTimestamp();
    sessionStorage.setItem(Varname.Settings_LoadTime, timeNow.toString());
}

function getSettings() {
    let settingsLoadTime = sessionStorage.getItem(Varname.Settings_LoadTime);
    if (settingsLoadTime == null) {
        console.error(Err.Settings);
        return null;
    }

    return new Settings(
        sessionStorage.getItem(Varname.Settings_Version),
        sessionStorage.getItem(Varname.Settings_TTL),
        sessionStorage.getItem(Varname.Settings_SiteName),
        sessionStorage.getItem(Varname.Settings_SiteDomain),
        sessionStorage.getItem(Varname.Settings_SessionMaxDuration),
        sessionStorage.getItem(Varname.Settings_MessageEditTime),
        sessionStorage.getItem(Varname.Settings_PageSize),
    );
}


// Entry point.
async function onPageLoad() {
    // Settings initialisation.
    let ok = await updateSettingsIfNeeded();
    if (!ok) {
        return;
    }
    let settings = getSettings();

    let curPage = window.location.search;

    drawPageHeader(settings);
    drawPageFooter(settings);

    //TODO
    let selfRoles = await getSelfRoles();
    if (selfRoles == null) {
        if (lastHttpStatusCode === httpStatusCode_NotAuthorised) {
            drawPageContent(settings, pageContent_logIn);
            return;
        }
    }

    console.log(selfRoles);
}

// UI functions.

function hideElement(el) {
    el.style.display = "none";
}

function showElement(el) {
    switch (el.tagName.toLowerCase()) {
        case "tr":
            el.style.display = "table-row";
            return;

        default:
            console.error("Unknown element type: " + el.tagName);
            return;
    }
}

function enableElement(el) {
    el.disabled = false;
}

function disableElement(el) {
    el.disabled = true;
}

function newDiv() {
    return document.createElement("DIV");
}

function newFieldset() {
    return document.createElement("FIELDSET");
}

function newTable() {
    return document.createElement("TABLE");
}

function newTr() {
    return document.createElement("TR");
}

function newTh() {
    return document.createElement("TH");
}

function newTd() {
    return document.createElement("TD");
}

function newInput() {
    return document.createElement("INPUT");
}

function drawPageHeader(settings) {
    let ph = document.getElementById(id_td_pageHeader);
    ph.textContent = settings.SiteName + " " + "header";
}

function drawPageFooter(settings) {
    let pf = document.getElementById(id_td_pageFooter);
    pf.textContent = settings.SiteName + " " + "footer";
}

function drawPageContent(settings, contentType) {
    let pc = document.getElementById(id_td_pageContent);

    switch (contentType) {
        case pageContent_logIn:
            drawPageContent_LogIn(settings, pc);
            return;

        default:
            console.error("Unknown page content type: " + contentType);
            return;
    }
}

function drawPageContent_LogIn(settings, pc) {
    pc.innerHTML = `
<table id="logIn">
    <tr>
        <td>E-Mail</td>
        <td>
            <input type="text" name="user_email"/>
        </td>
    </tr>
    <tr>
        <td colspan="2">
            <input type="button" name="log_in_proceed_1" value=" Proceed " onClick="on_log_in_proceed_1_click(this)"/>
        </td>
    </tr>
    <tr>
        <td>Captcha Question</td>
        <td>
            <img alt="captcha_question" src=""/>
        </td>
    </tr>
    <tr>
        <td>Captcha Answer</td>
        <td>
            <input type="text" name="captcha_answer"/>
        </td>
    </tr>
    <tr>
        <td>Verification Code</td>
        <td>
            <input type="text" name="verification_code"/>
        </td>
    </tr>
    <tr>
        <td>Password</td>
        <td>
            <input type="password" name="user_pwd"/>
        </td>
    </tr>
    <tr>
        <td colspan="2">
            <input type="button" name="log_in_proceed_2" value=" Proceed " onClick="on_log_in_proceed_2_click(this)"/>
        </td>
    </tr>
    <tr>
        <td>Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
    <tr>
        <td>Auth data</td>
        <td>
            <input type="text" name="auth_data"/>
        </td>
    </tr>
</table>`;
    let tbl = document.getElementById(id_table_logIn);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 1) {
            hideElement(tbl.rows[i]);
        }
    }
}

// Event handlers.

async function on_log_in_proceed_1_click(e) {
    let tbl = e.parentNode.parentNode.parentNode;
    let email = tbl.rows[0].children[1].children[0].value;
    console.log(email);
    let ok = validateEmailAddress(email);
    if (!ok) {
        console.error("E-Mail address is not valid");
        return;
    }

    let res = await startLogIn(email);
    if (res == null) {
        return;
    }

    // E-Mail address now can not be changed.
    disableElement(tbl.rows[0].children[1].children[0]);

    tbl.rows[7].children[1].children[0].value = res.requestId;
    tbl.rows[8].children[1].children[0].value = res.authData;
    let captchaId = res.captchaId;

    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 1) {
            showElement(tbl.rows[i]);
        }
    }
    hideElement(tbl.rows[1]);
    hideElement(tbl.rows[7]);
    disableElement(tbl.rows[7].children[1].children[0]);
    hideElement(tbl.rows[8]);
    disableElement(tbl.rows[8].children[1].children[0]);

    // Show captcha image.
    let captchaImg = tbl.rows[2].children[1].children[0];
    captchaImg.src = makeUrl_CaptchaImage(captchaId);
}

async function on_log_in_proceed_2_click(e) {
    let tbl = e.parentNode.parentNode.parentNode;
    let captchaAnswer = tbl.rows[3].children[1].children[0].value;
    let vCode = tbl.rows[4].children[1].children[0].value;
    let pwd = tbl.rows[5].children[1].children[0].value;
    let requestId = tbl.rows[7].children[1].children[0].value;
    let saltBA = base64ToByteArray(tbl.rows[8].children[1].children[0].value);

    if (captchaAnswer.length === 0) {
        console.error("Captcha answer is not set");
        return;
    }
    if (vCode.length === 0) {
        console.error("Verification code is not set");
        return;
    }
    if (pwd.length === 0) {
        console.error("Password is not set");
        return;
    }
    if (!isPasswordAllowed(pwd)) {
        console.error("Password is not allowed");
        return;
    }
    if (requestId.length === 0) {
        console.error("Request ID is not set");
        return;
    }

    let keyBA = makeHashKey(pwd, saltBA);
    let authChallengeResponse = byteArrayToBase64(keyBA);

    let res = await confirmLogIn(requestId, captchaAnswer, vCode, authChallengeResponse);
    if (res == null) {
        return;
    }

    await redirectPage(true, "");
}
