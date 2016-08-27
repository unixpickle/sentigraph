(function() {

  window.app = {};

  var CLASSIFIER_PATH = 'data/sentibayes.gz';
  var WEBWORKER_PATH = 'decompressor/decompressor.js';
  var XHR_DONE = 4;
  var HTTP_OK = 200;

  function initialize() {
    fetchDecompressed(function(err, data) {
      if (err) {
        errorLoading(err);
      } else {
        var classifier;
        try {
          var str = '';
          for (var i = 42, len = data.length; i < len; ++i) {
            str += String.fromCharCode(data[i]);
          }
          classifier = JSON.parse(str);
        } catch (e) {
          errorLoading(e);
          return;
        }
        classifierLoaded(new window.app.Classifier(classifier));
      }
    });
  }

  function errorLoading(err) {
    var loader = document.getElementById('loader');
    loader.textContent = 'Load failed: ' + err;
  }

  function classifierLoaded(classifier) {
    document.body.className = '';
    var rateButton = document.getElementById('rate-button');
    var rateText = document.getElementById('rate-text');
    var rating = document.getElementById('rating');
    rateButton.addEventListener('click', function() {
      var text = rateText.value;
      var score = classifier.classify(text);
      if (score > 0) {
        rating.className = 'positive';
      } else {
        rating.className = 'negative';
      }
      rating.textContent = 'Score ' + score.toFixed(3);
    });
  }

  function fetchData(callback) {
    var xhr = new XMLHttpRequest();
    xhr.responseType = "arraybuffer";
    xhr.open('GET', CLASSIFIER_PATH);
    xhr.send(null);

    xhr.onreadystatechange = function () {
      if (xhr.readyState === XHR_DONE) {
        if (xhr.status === HTTP_OK) {
          callback(null, new Uint8Array(xhr.response));
        } else {
          callback('status '+xhr.status, null);
        }
      }
    };
  }

  function fetchDecompressed(callback) {
    fetchData(function(err, compressed) {
      if (err) {
        callback(err, null);
        return;
      }
      w = new Worker(WEBWORKER_PATH);
      w.onmessage = function(e) {
        if (e.data[1]) {
          callback(e.data[1], null);
        } else {
          callback(null, new Uint8Array(e.data[0]));
        }
      };
      w.postMessage(compressed);
    });
  }

  window.addEventListener('load', initialize);

})();
