import { Feature } from "ol";
import { Point } from "ol/geom";
import Style from "ol/style/Style";
import OlCircle from "ol/style/Circle";
import OlFill from "ol/style/Fill";
import Stroke from "ol/style/Stroke";
import OlText from "ol/style/Text";
import { PlayerData } from "../api/player";
import Map from "ol/Map.js";
import VectorLayer from "ol/layer/Vector.js";
import VectorSource from "ol/source/Vector";
import { getGradientColor } from "../utils";
import { Assignment } from "../api/optimizer";
import Icon from "ol/style/Icon";

const makeHouse = (size: number) => {
  return new Icon({
    anchor: [0.5, 25],
    anchorXUnits: "fraction",
    anchorYUnits: "pixels",
    src: `/house_pop_${size}.png`,
  });
};

export const BASE_STYLE = new Style({
  image: makeHouse(48),
});

export const HOVER_STYLE = new Style({
  image: makeHouse(64),
});

export function updateBadgeStyles(map: Map, player: PlayerData | undefined) {
  // find the first vector layer
  const vectorLayer = map
    .getLayers()
    .getArray()
    .find((l) => l instanceof VectorLayer) as VectorLayer<
    VectorSource<Feature<Point>>
  >;

  if (!vectorLayer) return;

  vectorLayer
    .getSource()
    ?.getFeatures()
    .forEach((feature) => {
      feature.setStyle(createBadgeStyle(feature, player, getGradientColor));
    });
}

export function createBadgeStyle(
  feature: Feature<Point>,
  player: PlayerData | undefined,
  getGradientColor: (index: number) => string,
): Style[] {
  const pinId = (feature.get("id") as number) + 1;
  const prioritized = player?.plotData[pinId];

  if (prioritized === undefined) {
    return [BASE_STYLE];
  }

  const badgeStyle = new Style({
    image: new OlCircle({
      radius: 10,
      fill: new OlFill({ color: getGradientColor(prioritized) }),
      stroke: new Stroke({
        color: "#000000",
        width: 1,
      }),
      displacement: [-15, 0], // moves badge horizontally
    }),
    text: new OlText({
      text: prioritized.toString(),
      font: "bold 11px Geist, sans-serif",
      fill: new OlFill({ color: "black" }),
      textAlign: "center",
      textBaseline: "middle",
      offsetX: -15, // match the displacement
    }),
    zIndex: 1,
  });

  return [badgeStyle, BASE_STYLE];
}

export function createAssignmentBadge(map: Map, assignments: Assignment[]) {
  const vectorLayer = map
    .getLayers()
    .getArray()
    .find((l) => l instanceof VectorLayer) as VectorLayer<
    VectorSource<Feature<Point>>
  >;

  if (!vectorLayer) return;

  vectorLayer
    .getSource()
    ?.getFeatures()
    .forEach((feature) => {
      feature.setStyle(
        createAssignmentStyle(feature, assignments, getGradientColor),
      );
    });
}

export function createAssignmentStyle(
  feature: Feature<Point>,
  assignments: Assignment[],
  getGradientColor: (index: number) => string,
): Style[] {
  const pinId = (feature.get("id") as number) + 1;
  const ass = assignments.find((ass) => ass.plot === pinId);

  if (ass === undefined) {
    return [BASE_STYLE];
  }

  const badgeStyle = new Style({
    text: new OlText({
      text: ass.char.toString(), // longer text allowed
      font: "bold 14px Geist, sans-serif",
      fill: new OlFill({ color: "black" }),
      stroke: new Stroke({ color: getGradientColor(ass.score), width: 2 }), // optional outline for readability
      textAlign: "center",
      textBaseline: "middle",
      offsetY: -15,
      padding: [2, 6, 2, 6],
    }),
    zIndex: 1,
  });

  return [badgeStyle, BASE_STYLE];
}
