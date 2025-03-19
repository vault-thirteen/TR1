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

    //TODO
    //let x = await addForum("Forum Y");
    //console.log(x);
}
