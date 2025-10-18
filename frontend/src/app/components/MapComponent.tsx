"use client";

import { useEffect, useRef, useState } from "react";

// OpenLayers imports: Using explicit .js file extensions for bundler compatibility
import Feature from "ol/Feature.js";
import Map from "ol/Map.js";
import View from "ol/View.js";
import Point from "ol/geom/Point.js";
import ImageLayer from "ol/layer/Image.js";
import VectorLayer from "ol/layer/Vector.js";
import VectorSource from "ol/source/Vector.js";
import Style from "ol/style/Style.js";
// We only need the type imports for the rest of the styles/events
import { getCenter } from "ol/extent";
import "ol/ol.css";
import { Projection } from "ol/proj.js";
import Static from "ol/source/ImageStatic.js";
import { PlotData, plotData } from "../data/PlotData";

import "@/styles/MapStyles.css";
import Icon from "ol/style/Icon";
import { getCommunityData, PlayerData } from "../api/player";
import { useAuth } from "../context/AuthContext";
import { getLowestFreePriority } from "../utils";
import ControlPanel from "./ControlPanel";
import { createAssignmentBadge, updateBadgeStyles } from "./Features";
import MapHoverPopup from "./Popup";
import { Assignment, getAssignedPlots } from "../api/optimizer";

/**
 * OpenLayers Map Component for displaying a static image with clickable pins.
 * The map uses an ImageStatic source, treating the image's dimensions as its coordinate system.
 */
export default function MapComponent() {
  const mapRef = useRef(null);
  const mapInstanceRef = useRef<Map>(null);

  const { user } = useAuth();
  const lockedRef = useRef(user?.community.locked ?? false);

  const [playerData, setPlayerData] = useState<PlayerData[]>([]);
  const [plotAssignments, setPlotAssignments] = useState<Assignment[]>([]);
  const [loading, setLoading] = useState<boolean>(true);

  const playerRef = useRef<PlayerData | undefined>(undefined);
  const player = playerData?.find(
    (player) => player.battletag === user?.battletag,
  );

  const rerenderFeatures = () => {
    if (!mapInstanceRef.current) return;
    if (plotAssignments?.length > 0) {
      createAssignmentBadge(mapInstanceRef.current, plotAssignments, baseStyle);
    } else {
      updateBadgeStyles(mapInstanceRef.current, playerRef.current, baseStyle);
    }
  };

  useEffect(() => {
    playerRef.current = player;
    lockedRef.current = user?.community.locked ?? false;
    rerenderFeatures();
  }, [user?.community.locked, playerData, plotAssignments]);

  const updatePlayerPlot = (plotId: number, value: number) => {
    if (user?.community.locked) return;
    setPlayerData((prev) =>
      prev.map((p) =>
        p.battletag === user?.battletag
          ? {
              ...p,
              plotData: { ...p.plotData, [plotId]: value },
            }
          : p,
      ),
    );
  };

  useEffect(() => {
    async function fetchData() {
      try {
        if (user?.community.locked) {
          const data = await getAssignedPlots();
          setPlotAssignments(data);
        }
        const data = await getCommunityData();
        setPlayerData(data);
        playerRef.current = data?.find(
          (player) => player.battletag === user?.battletag,
        );
      } catch (err: any) {
        console.error(err);
      } finally {
        setLoading(false);
      }
    }
    fetchData();
  }, []);

  const imageUrl = "/housing_map.jpg";
  const imageExtent = [0, 0, 3840, 2560];
  const baseStyle = new Style({
    image: new Icon({
      anchor: [0.5, 25],
      anchorXUnits: "fraction",
      anchorYUnits: "pixels",
      src: "/house_pop_48.png",
    }),
  });

  useEffect(() => {
    if (
      !mapRef.current ||
      mapInstanceRef.current ||
      (!playerRef.current && plotAssignments.length == 0)
    ) {
      return;
    }

    const vectorSource = new VectorSource({
      features: plotData.map((plot, index) => {
        const feature = new Feature({
          geometry: new Point([plot.xCoord, plot.yCoord]),
          name: plot.label,
          id: index,
          plot: plot,
        });
        feature.setStyle(baseStyle);
        return feature;
      }),
    });

    const vectorLayer = new VectorLayer({
      source: vectorSource,
      extent: imageExtent,
    });

    const projection = new Projection({
      code: "housing-map",
      units: "pixels",
      extent: imageExtent,
    });

    const imageLayer = new ImageLayer({
      source: new Static({
        url: imageUrl,
        projection: projection,
        imageExtent: imageExtent,
      }),
    });

    const map = new Map({
      layers: [imageLayer, vectorLayer],
      target: mapRef.current,
      view: new View({
        projection: projection,
        center: getCenter(imageExtent),
        zoom: 3.2,
        maxZoom: 6,
      }),
    });

    mapInstanceRef.current = map;

    if (plotAssignments?.length > 0) {
      createAssignmentBadge(mapInstanceRef.current, plotAssignments, baseStyle);
    } else {
      updateBadgeStyles(mapInstanceRef.current, playerRef.current, baseStyle);
    }

    map.on("click", function (evt) {
      map.forEachFeatureAtPixel(evt.pixel, function (feature, layer) {
        if (
          feature &&
          feature.getGeometry()?.getType() === "Point" &&
          !lockedRef.current
        ) {
          const plot = feature.get("plot") as PlotData;
          updatePlayerPlot(plot.id, getLowestFreePriority(playerRef.current!));
        }
      });
    });
  }, [playerData, plotAssignments]);

  return (
    <div className="component-style">
      <div className="map-wrapper">
        {/* Pin Info Panel */}
        <ControlPanel
          user={user}
          playerData={player}
          updatePlayerPlot={updatePlayerPlot}
          updatePlotAssignments={setPlotAssignments}
        />
        {/* Map Container */}
        <div ref={mapRef} id="map" className="map-container-style" />
        {mapInstanceRef.current && (
          <MapHoverPopup map={mapInstanceRef.current} />
        )}
      </div>
    </div>
  );
}
