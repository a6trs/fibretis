sight_items = [];
sight_dropdown = document.getElementById('sight-btn');
sight_lastsel = 0;

disp_sel_sight_item = function (idx) {
  'use strict';
  sight_dropdown.innerHTML = sight_items[idx].innerHTML + "&nbsp;<span class='am-icon-caret-down'></span>";
  sight_items[sight_lastsel].classList.remove('am-active');
  sight_items[idx].classList.add('am-active');
  sight_lastsel = idx;
};

init_sight_dropdown = function (prjid) {
  'use strict';

  var sight_selctd = function (idx) {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/sight/projects/' + prjid + '/' + idx);
    xhr.send(null);
    disp_sel_sight_item(idx);
  };

  var handler = function (idx) {
    return function () { sight_selctd(idx); };
  };

  var i;
  for (i = 0; i < 3; i++) {
    sight_items[i] = document.getElementById('sight-sel-' + i.toString());
    sight_items[i].onclick = handler(i);
  }

};
