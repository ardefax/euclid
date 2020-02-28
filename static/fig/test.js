const svgNS = 'http://www.w3.org/2000/svg';

function N(svgAnimatedLength) { //TODO
  return svgAnimatedLength.baseVal.value;
}
function T(elem) {
  return elem.classList[0] // TODO also elem.tagName?
}

class Ui {
  constructor(svg) {
    this.svg = svg;
    this.handle = null;

    svg.addEventListener("pointerdown", evt => this.pointerDown(evt));
    svg.addEventListener("pointermove", evt => this.pointerMove(evt));
    svg.addEventListener("pointerup", evt => this.pointerUp(evt));
  }

  pointerDown(evt) {
    const [x, y] = this.getCoords(evt)
    this.handle = null;
    if (evt.target.classList.contains("point")) {
      this.handle = evt.target;
    }
    console.log('pointer:down x', x, 'y', y, 'h', this.handle);
    if (this.handle == null) {
      return
    }

    // TODO Ideally move the pointer correctly relative to the mouse
    const [hx, hy] = this.origin(this.handle),
          dx = x - hx,
          dy = y - hy;
    this.move(this.handle, dx, dy);
    this.update(this.handle, dx, dy);
  }
  pointerMove(evt) {
    const [x, y] = this.getCoords(evt);
    //console.log('pointer:move x', x, 'y', y, 'h', this.handle);
    if (this.handle == null) {
      return;
    }

    const [hx, hy] = this.origin(this.handle),
          dx = x - hx,
          dy = y - hy;
    this.move(this.handle, dx, dy);
    this.update(this.handle, dx, dy);
  }
  pointerUp(evt) {
    const [x, y] = this.getCoords(evt)
    console.log('pointer:up x', x, 'y', y, 'h', this.handle);
    if (this.handle == null) {
      return;
    }

    const [hx, hy] = this.origin(this.handle),
          dx = x - hx,
          dy = y - hy;
    this.move(this.handle, dx, dy);
    this.update(this.handle, dx, dy);
    this.handle = null;
  }

  // https://github.com/raphlinus/spline-research/blob/1a0fd3df09db517726309e899f165dc225e466e3/splineui.js#L662-L670
	// On Chrome, just offsetX, offsetY work, but on FF it takes the group transforms
	// into account. We always want coords relative to the SVG.
	getCoords(evt) {
    const CTM = this.svg.getScreenCTM();
		let rect = this.svg.getBoundingClientRect();
    //console.log("svg:rect:", rect);
    //console.log("evt:client:", evt.clientX, evt.clientY);
		let x = (evt.clientX - rect.left - CTM.e) / CTM.a;
		let y = (evt.clientY - rect.top - CTM.f) / CTM.d;
		return [x, y];
	}

  update(src, dx, dy) {
    const svg = this.svg,
          processed = {};
    processed[src.id] = src;

    let cascade = Array.from(svg.getElementsByClassName(src.id));
    while (cascade.length > 0) {
      const elem = cascade.shift();
      //if (!!processed[elem.id]) {
      //  continue;
      //}
      //processed[elem.id] = elem;

      this.redraw(elem);
      cascade = cascade.concat(Array.from(svg.getElementsByClassName(elem.id)));
      // TODO Recursively collect what needs updating...
      // TODO Validate that my "class pointers" form a DAG
    }
  }

  redraw(elem) {
    const svg = this.svg,
          tag = elem.tagName,
          def = elem.classList;
    switch (def[0]) { // TODO also care about tag?
      case "circle": { // class='circle A B'
        const center = svg.getElementById(def[1]),
              radius = svg.getElementById(def[2]),
              [cx, cy] = this.origin(center),
              [rx, ry] = this.origin(radius);

        console.log("redraw: circle", cx, cy, rx, ry);

        elem.setAttribute('cx', cx);
        elem.setAttribute('cy', cy);
        elem.setAttribute('r', Math.hypot(cx-rx, cy-ry));
      } break;
      case "line": {
        const p1 = svg.getElementById(def[1]),
              p2 = svg.getElementById(def[2]),
              [x1, y1] = this.origin(p1),
              [x2, y2] = this.origin(p2);

        elem.setAttribute('x1', x1);
        elem.setAttribute('y1', y1);
        elem.setAttribute('x2', x2);
        elem.setAttribute('y2', y2);
      } break;
      case "intersection": {
        const p1 = svg.getElementById(def[1]),
              p2 = svg.getElementById(def[2]),
              [x, y] = this.intersect(p1, p2);
        // TODO validate that tag is a circle
        elem.setAttribute('cx', x);
        elem.setAttribute('cy', y);
      } break
      default:
        console.warn("redraw: TODO tag:def", tag, def);
    }
  }

  intersect(elem1, elem2) {
    const t1 = T(elem1), t2 = T(elem2);
    if (t1 == "circle" && t2 == "circle") {
      // https://stackoverflow.com/a/3349134
      const [cx1, cy1] = this.origin(elem1),
            [cx2, cy2] = this.origin(elem2),
            r1 = N(elem1.r),
            r2 = N(elem2.r),
            d = Math.hypot(cx2-cx1, cy2-cy1);
      if (d > (r1 + r2) || d < Math.abs(r2 - r1)) {
        return [];
      }
      // TODO d = 0 and r1 == r2 coincident means infinite points

      // r1^2 = h^2 + a^2; r2^2 = h^2 + b^2
      //    where d = a + b
      const a = (r1*r1 - r2*r2 + d*d) / (2*d),
            h = Math.sqrt(r1*r1 - a*a),
            cdx = (cx2 - cx1) / d,
            cdy = (cy2 - cy1) / d,
            lx = cx1 + a*cdx, // roughly lens midpoint
            ly = cy1 + a*cdy;

      //console.log({a, h, lx, ly, cdx, cdy});
      return [
        lx + h*cdy, ly - h*cdx,
        lx - h*cdy, ly + h*cdx
      ];
    }
    console.warn("redraw: TODO tag:def", tag, def);
  }

  nudge(src, elem, dx, dy) {
    //console.log("nudge:", src, elem)
    const tag = elem.tagName;
    switch (tag) {
      case "circle":
        const cx = N(elem.cx),
              cy = N(elem.cy);
        // classList is `center A radius B`
        if (src == elem.id || elem.classList[1] == src) {
          // Moving the center point
          elem.setAttribute('cx', cx + dx);
          elem.setAttribute('cy', cx + dy);
        } else {
          // Moving the radius point
          const [ox, oy] = this.origin(elem),
                r = N(elem.r),
                x = ox + r + dx, y = oy + dy;
          //debugger;
          console.log("nudge:radius o:", ox, oy, "r:", r, "x", x, "y", y);
          elem.setAttribute('r', Math.hypot(ox, x, oy, y));
        }
        break;
      case "line":
        if (src == elem.classList[1]) {
          elem.setAttribute('x1', N(elem.x1) + dx);
          elem.setAttribute('y1', N(elem.y1) + dy);
        } else {
          elem.setAttribute('x2', N(elem.x2) + dx);
          elem.setAttribute('y2', N(elem.y2) + dy);
        }
        break;
      default:
        console.warn("nudge: missing-tag", tag, elem);
    }
  }
  origin(elem) { // TODO anchor? pivot?
    const tag = elem.tagName;
    switch (tag) {
      case "circle":
      case "ellipse":
        return [N(elem.cx), N(elem.cy)];
      case "line":
        return [N(elem.x1), N(elem.y1)];
      default:
        console.warn("origin: missing-tag", tag, elem)
    }
    return [0, 0]; // TODO Better "wrong" values
  }
  move(elem, dx, dy) {
    const tag = elem.tagName;
    switch (tag) {
      case "circle":
      case "ellipse":
        elem.setAttribute('cx', N(elem.cx) + dx);
        elem.setAttribute('cy', N(elem.cy) + dy);
        break;
      case "line":
        elem.setAttribute('x1', N(elem.x1) + dx);
        elem.setAttribute('y1', N(elem.y1) + dy);
        elem.setAttribute('x2', N(elem.x2) + dx);
        elem.setAttribute('y2', N(elem.y2) + dy);
        break;
      default:
        console.warn("move: missing-tag", tag, elem);
    }
  }
}

window.setTimeout(() => { // XXX Hack for embedded svg object to be loaded.
  const objects = document.getElementsByTagName('object')
  const svg = objects[0].contentDocument.firstElementChild;
  const ui = new Ui(svg);
  console.log("new UI");
}, 1000);
