import React, { useEffect, useRef } from "react";
import Overlay from "ol/Overlay";
import { Map } from "ol";
import "ol/ol.css";
import { Point } from "ol/geom";

interface MapHoverPopupProps {
  map?: Map;
}

export default function MapHoverPopups({ map }: MapHoverPopupProps) {
  const popupRef = useRef<HTMLDivElement | null>(null);
  const overlayRef = useRef<Overlay | null>(null);

  useEffect(() => {
    if (!map || !popupRef.current) return;

    overlayRef.current = new Overlay({
      element: popupRef.current,
      positioning: "bottom-center",
      stopEvent: false,
      offset: [0, -25],
    });

    map.addOverlay(overlayRef.current);

    const handlePointerMove = (evt: any) => {
      const feature = map.forEachFeatureAtPixel(evt.pixel, (feat) => feat);
      const mapEl = map.getTargetElement() as HTMLElement;
      mapEl.style.cursor = feature ? "pointer" : "";

      if (feature) {
        const pointGeometry = feature.getGeometry() as Point;
        const coords = pointGeometry.getCoordinates();
        const pinId = feature.get("plot").id;

        if (overlayRef.current && coords) {
          popupRef.current!.innerHTML = `Plot #${pinId}`;
          overlayRef.current.setPosition(coords);
        }
      } else {
        overlayRef.current?.setPosition(undefined);
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
