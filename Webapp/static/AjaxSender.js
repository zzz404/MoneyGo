function AjaxSender() {
    if (this.constructor !== AjaxSender) {
        return new AjaxSender();
    }
    this.xhr = new XMLHttpRequest();
}

(function () {
    AjaxSender.prototype.httpGet = function (url, onSuccess, onFail) {
        xhr = prepareXhr(this, url, "GET", onSuccess, onFail);
        xhr.send();
    };

    function prepareXhr(sender, url, method, onSuccess, onFail) {
        var xhr = sender.xhr;
        xhr.open(method, url, true);
        xhr.onload = function () {
            if (xhr.status === 200) {
                onSuccess(sender);
            } else {
                if (onFail) {
                    onFail(sender);
                } else {
                    alert("ajax fail : " + xhr.status);
                }
            }
        };
        return xhr;
    }

    AjaxSender.prototype.httpPost = function (url, data, onSuccess, onFail) {
        var ps = [];
        for (p in data) {
            ps.push(p + "=" + data[p]);
        }
        var dataString = ps.join("&");

        xhr = prepareXhr(this, url, "POST", onSuccess, onFail);
        xhr.setRequestHeader(
            "Content-Type",
            "application/x-www-form-urlencoded"
        );
        xhr.send(encodeURI(dataString));
    };

    AjaxSender.prototype.checkSuccess = function () {
        var result = JSON.parse(this.xhr.responseText);
        checkSuccess(result);
    };

    function checkSuccess(result) {
        if (!result.Success) {
            var s = "Error: " + result.Message;
            alert(s);
            throw s;
        }
    }

    AjaxSender.prototype.getJsonData = function () {
        var result = JSON.parse(this.xhr.responseText);
        checkSuccess(result);
        return result.data;
    };
})();
