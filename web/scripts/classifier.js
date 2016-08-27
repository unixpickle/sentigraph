(function() {

  function Classifier(jsObj) {
    this._negative = jsObj.Conditional['1'];
    this._positive = jsObj.Conditional['2'];
  }

  Classifier.prototype.classify = function(sampleText) {
    var sampleKeywords = keywordFlags(sampleText);
    var negProb = 0, posProb = 0;

    var keywords = Object.keys(this._negative);
    for (var i = 0, len = keywords.length; i < len; ++i) {
      var keyword = keywords[i];
      if (sampleKeywords[keyword]) {
        negProb += Math.log(this._negative[keyword]);
        posProb += Math.log(this._positive[keyword]);
      } else {
        negProb += Math.log(1 - this._negative[keyword]);
        posProb += Math.log(1 - this._positive[keyword]);
      }
    }

    return posProb - negProb;
  };

  window.app.Classifier = Classifier;

  function keywordFlags(text) {
    var map = {};
    var keywords = normalizedKeywords(text);
    for (var i = 0, len = keywords.length; i < len; ++i) {
      map[keywords[i]] = true;
    }
    return map;
  }

  function normalizedKeywords(text) {
    var res = [];
    var tokens = splitText(text);
    for (var i = 0, len = tokens.length; i < len; ++i) {
      var tok = tokens[i].toLowerCase();
      if (tok[0] === '@') {
        res.push('USERNAME');
      } else if (tok.match(/https?:\/\//)) {
        res.push('URL');
      } else {
        res.push.apply(res, separatePunctuation(removeRepeatedLetters(tok)));
      }
    }
    return res;
  }

  function splitText(text) {
    return text.split(/\s/);
  }

  function removeRepeatedLetters(text) {
    var last = '';
    var count = 0;
    var res = '';
    for (var i = 0, len = text.length; i < len; ++i) {
      var ch = text[i];
      if (last == ch) {
        ++count;
      } else {
        last = ch;
        count = 1;
      }
      if (count <= 2) {
        res += ch;
      }
    }
    return res;
  }

  function separatePunctuation(text) {
    punctuation = ['!', '.', ',', '?', '"', '(', ')'];

    var words = [];
    var cur = '';
    var last = '';

    for (var i = 0, len = text.length; i < len; ++i) {
      var ch = text[i];
      if (!cur) {
        last = ch;
        cur += ch;
        continue;
      }
      var isPunct = (punctuation.indexOf(ch) >= 0);
      var wasPunct = (punctuation.indexOf(last) >= 0);
      if (isPunct != wasPunct) {
        words.push(cur);
        cur = ch;
      } else {
        cur += ch;
      }
      last = ch;
    }

    if (cur) {
      words.push(cur);
    }
    return words;
  }

})();
