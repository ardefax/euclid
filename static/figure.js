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

    // Build the DAG for updates
  }

  pointerDown(evt) {
    const [x, y] = this.getCoords(evt)
    this.handle = null;
    // Need larger hit tests
    // - One trick would be mimicking a touch event that expands beyound the hit target
    if (evt.target.classList.contains("point")) {
      this.handle = evt.target;
    }
    console.log('pointer:down x', x, 'y', y, 'h', this.handle);
    var req = new XMLHttpRequest();
    req.open("GET", "/down?x=" + x + "&y=" + y + "&type=" + evt.pointerType);
    req.send();
    if (this.handle == null) {
      return
    }
    //this.handle.setPointerCapture(evt.pointerId);

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
    var req = new XMLHttpRequest();
    req.open("GET", "/up?x=" + x, "&y=", y);
    req.send();
    if (this.handle == null) {
      return;
    }
    //this.handle.releasePointerCapture(evt.pointerId);

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
    console.log("------ update --------")

    let counter = 0;

    let cascade = Array.from(svg.getElementsByClassName(src.id));
    while (cascade.length > 0) {
      counter++;
      if (counter > 100) {
        debugger; // TOO Blowing the stack
        break;
      }

      const elem = cascade.shift();
      if (!!processed[elem.id]) {
        console.log("skipping", elem.tagName, elem.id)
        continue;
      }

      // Ensure we've already processed the dependents of this list
      const requeue = Array.from(elem.classList).some((depId) => {
        // TODO This isn't quite right either since it's going to find
        // things that may not be part of the DAG rooted at the original
        // `src` that changed to cause the update. Really need to just
        // build the DAG and do a proper depth-first-search update.
        const dep = svg.getElementById(depId)
        if (!!dep && !processed[depId]) {
          cascade.push(dep);
          return true;
        }
        return false;
      }, []);
      if (requeue) {
        // Re-queue this element since there are missing deps
        // TODO This could cycle forever if we don't have a DAG
        console.log("requeuing", elem.id);
        cascade.push(elem);
        continue;
      }

      this.redraw(elem);
      cascade = cascade.concat(Array.from(svg.getElementsByClassName(elem.id)));
      processed[elem.id] = elem;
    }
  }

  redraw(elem) {
    const svg = this.svg,
          tag = elem.tagName,
          def = elem.classList;

    console.log("redraw", elem)
    switch (def[0]) { // TODO also care about tag?
      case "point": {
        // XXX No-op since these are to trigger the initial draws from input
        // TODO Could do something with "random" points here...
      } break;

      case "circle": { // class='circle A B'
        const center = svg.getElementById(def[1]),
              radius = svg.getElementById(def[2]),
              [cx, cy] = this.origin(center),
              [rx, ry] = this.origin(radius);

        //console.log("redraw: circle", cx, cy, rx, ry);

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
      } break;

      case "label": { // label A left AB; label B right AB; label C above AC BC
        // label <target> [constraint] TODO Multiple constraints
        const bbox = elem.getBBox(),
          constraints = Array.from(def).slice(2),
          direction = constraints[0],
          params = constraints.slice(1);

        elem.setAttribute('text-anchor', 'middle');
        elem.setAttribute('dominant-baseline', 'central');

        switch (direction) {
          case "left": { // relative to x1,y1
            const [, x1, y1, x2, y2, d] = this.lineAnchor(params[0]);
            elem.setAttribute('x', x1 + bbox.width * (x1 - x2) / d);
            elem.setAttribute('y', y1 + bbox.height * (y1 - y2) / d);
          } break;
          case "right": { // relative to x2,y2
            const [, x1, y1, x2, y2, d] = this.lineAnchor(params[0]);
            elem.setAttribute('x', x2 + bbox.width * (x2 - x1) / d);
            elem.setAttribute('y', y2 + bbox.height * (y2 - y1) / d);
          } break;
          case "above": { // "between" the rays for the two anchor, e.g.
            // the vector that bisects the two lines. Assumes that the
            // two lines share the same x2,y2 points
            const [, x11, y11, x12, y12, d1] = this.lineAnchor(params[0]);
            const [, x21, y21, x22, y22, d2] = this.lineAnchor(params[1]);
            const v1 = [(x12 - x11) / d1, (y12 - y11) / d1];
            const v2 = [(x22 - x21) / d2, (y22 - y21) / d2];
            const v = [ v1[0]+v2[0], v1[1]+v2[1] ];
            const d = Math.hypot(v[0], v[1]);

            if (x12 != x22 || y12 != y22) {
              consale.warn(`redraw: above ${params[0]} ${params[1]} don't share endpoint: (${x12},${y12}) != (${x22},${y22})`)
            }
            elem.setAttribute('x', x12 + bbox.width * v[0] / d);
            elem.setAttribute('y', y12 + bbox.height * v[1] / d);
          } break;
          case "incirc": { // incirc ACE <rads>  TODO this name sucks
            // where dregrees is CCW relative to the radial point used to describe the circle
            const [, cx, cy, r, px, py] = this.circleAnchor(params[0]),
              rads = params[1],
              sin = Math.sin(rads),
              cos = Math.cos(rads),
              dx = px - cx,
              dy = py - cy;

            // Formula for vector rotation in 2D
            // https://matthew-brett.github.io/teaching/rotation_2d.html
            elem.setAttribute('x', cx + dx * cos - dy * sin);
            elem.setAttribute('y', cy + dx * sin + dy * cos);
          } break;
          default:
            console.warn("redraw: unexpected direction tag:def", tag, def)
        }
        
      } break;

      default:
        console.warn("redraw: TODO tag:def", tag, def);
    }
  }

  lineAnchor(id) {
      const anchor = this.svg.getElementById(id);
      if (!anchor) {
        console.warn("redraw: missing anchor id:", id)
        return [];
      }
      const [type, x1, y1, x2, y2] = this.decompose(anchor);
      if (type != 'line') {
        console.warn("redraw: anchor not a line id:type", id, type)
        return [type];
      }
      return [type, x1, y1, x2, y2, Math.hypot(x1-x2, y1-y2)];
  }
  circleAnchor(id) {
      const anchor = this.svg.getElementById(id);
      if (!anchor) {
        console.warn("redraw: missing anchor id:", id)
        return [];
      }
      const [type, cx, cy, r] = this.decompose(anchor);
      if (type != 'circle') {
        console.warn("redraw: anchor not a circle id:type", id, type)
        return [type];
      }
      // Resolve the other point on the circle as well, e.g class='circle <center> <radius-point>'
      const [,, px, py] = this.decomposeId(anchor.classList[2])
      return [type, cx, cy, r, px, py];
  }

  decomposeId(id) {
    const elem = this.svg.getElementById(id);
    if (!elem) {
      console.warn("decomposie: missing elem id:", id)
      return [];
    }
    const arr = this.decompose(elem);
    return [elem, ...arr];
  }
  decompose(elem) {
    const t = T(elem);
    switch (t) {
      case 'point': // <circle> under-the-hood
        return [t, N(elem.cx), N(elem.cy)];
      case 'line':
        return [t, N(elem.x1), N(elem.y1), N(elem.x2), N(elem.y2)];
      case 'circle':
        return [t, N(elem.cx), N(elem.cy), N(elem.r)];
      case 'ellipse':
        return [t, N(elem.cx), N(elem.cy), N(elem.r)]; // TODO r1 and r2
      default:
        console.warn("decompose: unexpected", t);
    }
    return [];
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

const object = document.getElementById('object-svg');
object.addEventListener('load', function(evt) {
  const svg = object.contentDocument.firstElementChild;
  const ui = new Ui(svg);
  ui.update({ id: "point" });
}, false);
