import React, { useEffect, useRef } from "react";
import Overlay from "ol/Overlay";
import { Map } from "ol";
import "ol/ol.css";
import { Point } from "ol/geom";
import Feature from "ol/Feature.js";
import { BASE_STYLE, HOVER_STYLE } from "./Features";
import Icon from "ol/style/Icon";
import Style from "ol/style/Style";

interface MapHoverPopupProps {
  map?: Map;
}

export default function MapHoverPopups({ map }: MapHoverPopupProps) {
  const popupRef = useRef<HTMLDivElement | null>(null);
  const overlayRef = useRef<Overlay | null>(null);

  const setFeatureImage = (feature: Feature<Point>, style: Style) => {
    const currentStyles = feature.getStyle();
    if (Array.isArray(currentStyles)) {
      const scaled = currentStyles.map((s) => {
        const img = s.getImage?.();
        if (img instanceof Icon) {
          return style;
        }
        return s;
      });
      feature.setStyle(scaled);
    }
  };

  useEffect(() => {
    if (!map || !popupRef.current) return;

    overlayRef.current = new Overlay({
      element: popupRef.current,
      positioning: "bottom-center",
      stopEvent: false,
      offset: [0, -25],
    });

    map.addOverlay(overlayRef.current);

    let lastFeature: any = null;

    const handlePointerMove = (evt: any) => {
      const feature =
        map.forEachFeatureAtPixel(
          evt.pixel,
          (feat): Feature<Point> | undefined => {
            if (
              feat instanceof Feature &&
              feat.getGeometry() instanceof Point
            ) {
              return feat as Feature<Point>;
            }
            return undefined;
          },
        ) ?? null;

      const mapEl = map.getTargetElement() as HTMLElement;
      mapEl.style.cursor = feature ? "pointer" : "";

      if (feature) {
        const pointGeometry = feature.getGeometry() as Point;
        const coords = pointGeometry.getCoordinates();
        const pinId = feature.get("plot").id;

        if (feature !== lastFeature) {
          if (feature) {
            setFeatureImage(feature, HOVER_STYLE);
          }
          if (lastFeature) {
            setFeatureImage(lastFeature, BASE_STYLE);
          }
          lastFeature = feature;
        }

        if (overlayRef.current && coords) {
          popupRef.current!.innerHTML = `Plot #${pinId}`;
          overlayRef.current.setPosition(coords);
        }
      } else {
        overlayRef.current?.setPosition(undefined);
        if (lastFeature) {
          setFeatureImage(lastFeature, BASE_STYLE);
          lastFeature = null;
        }
      }
    };

    map.on("pointermove", handlePointerMove);

    return () => {
      map.un("pointermove", handlePointerMove);
      if (overlayRef.current) map.removeOverlay(overlayRef.current);
    };
  }, [map]);

  return (
    <div
      ref={popupRef}
      style={{
        position: "absolute",
        background: "rgba(0, 0, 0, 0.75)",
        color: "#fff",
        padding: "4px 8px",
        borderRadius: "6px",
        fontSize: "12px",
        whiteSpace: "nowrap",
        pointerEvents: "none",
        transform: "translate(-50%, -100%)",
      }}
    />
  );
}
