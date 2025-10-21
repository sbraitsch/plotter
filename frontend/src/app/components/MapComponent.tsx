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
import { getCenter } from "ol/extent";
import "ol/ol.css";
import { Projection } from "ol/proj.js";
import Static from "ol/source/ImageStatic.js";
import { PlotData, plotData } from "../data/PlotData";

import "@/styles/MapStyles.css";
import { getCommunityData, PlayerData } from "../api/player";
import { useAuth } from "../context/AuthContext";
import { getLowestFreePriority } from "../utils";
import ControlPanel from "./ControlPanel";
import {
  BASE_STYLE,
  createAssignmentBadge,
  updateBadgeStyles,
} from "./Features";
import MapHoverPopup from "./Popup";
import { Assignment, getAssignedPlots } from "../api/optimizer";
import TargetedModal from "./TargetedModal";

/**
 * OpenLayers Map Component for displaying a static image with clickable pins.
 * The map uses an ImageStatic source, treating the image's dimensions as its coordinate system.
 */
export default function MapComponent() {
  const mapRef = useRef(null);
  const mapInstanceRef = useRef<Map>(null);

  const [contextDirty, setContextDirty] = useState(false);

  const [isTargetedModalOpen, setIsTargetedModalOpen] = useState(false);
  const [selectedPlot, setSelectedPlot] = useState<number | undefined>(
    undefined,
  );
  const handleOpenModal = (plotId: number) => {
    setSelectedPlot(plotId);
    setIsTargetedModalOpen(true);
  };

  const handleModalSubmit = async (plot: number, prio: number) => {
    forcePlotUpdate(plot, prio);
    setSelectedPlot(undefined);
    setIsTargetedModalOpen(false);
  };

  const { user } = useAuth();
  const lockedRef = useRef(user?.community.locked ?? false);

  const [targetedMode, setTargetedMode] = useState(false);
  const targetedRef = useRef(targetedMode);

  const [playerData, setPlayerData] = useState<PlayerData[]>([]);
  const [plotAssignments, setPlotAssignments] = useState<Assignment[]>([]);
  const [mapReady, setMapReady] = useState<boolean>(false);

  const playerRef = useRef<PlayerData | undefined>(undefined);
  const player = playerData?.find(
    (player) => player.battletag === user?.battletag,
  );

  const rerenderFeatures = () => {
    if (!mapInstanceRef.current) return;
    if (plotAssignments?.length > 0) {
      createAssignmentBadge(mapInstanceRef.current, plotAssignments);
    } else {
      updateBadgeStyles(mapInstanceRef.current, playerRef.current);
    }
  };

  useEffect(() => {
    targetedRef.current = targetedMode;
  }, [targetedMode]);

  useEffect(() => {
    playerRef.current = player;
    lockedRef.current = user?.community.locked ?? false;
    rerenderFeatures();
  }, [user?.community.locked, playerData, plotAssignments]);

  const clearPlayerMappings = () => {
    setPlayerData((prev) =>
      prev.map((p) =>
        p.battletag === user?.battletag
          ? {
              ...p,
              plotData: [],
            }
          : p,
      ),
    );
    setContextDirty(true);
  };

  const forcePlotUpdate = (plotId: number, value: number) => {
    if (user?.community.locked) return;
    setPlayerData((prev) =>
      prev.map((p) => {
        if (p.battletag !== user?.battletag) return p;

        const newPlotData = Object.fromEntries(
          Object.entries(p.plotData).filter(([_, v]) => v !== value),
        );

        newPlotData[plotId] = value;

        return {
          ...p,
          plotData: newPlotData,
        };
      }),
    );
    setContextDirty(true);
  };

  const updatePlayerPlot = (plotId: number, value: number) => {
    if (user?.community.locked) return;
    if (playerRef.current?.plotData[plotId]) {
      setPlayerData((prev) =>
        prev.map((p) =>
          p.battletag === user?.battletag
            ? {
                ...p,
                plotData: Object.fromEntries(
                  Object.entries(p.plotData).filter(
                    ([id]) => Number(id) !== plotId,
                  ),
                ),
              }
            : p,
        ),
      );
    } else {
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
    }
    setContextDirty(true);
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
      }
    }
    fetchData();
  }, []);

  const imageUrl = "/housing_map.jpg";
  const imageExtent = [0, 0, 3840, 2560];

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
        feature.setStyle(BASE_STYLE);
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
    setMapReady(true);

    if (plotAssignments?.length > 0) {
      createAssignmentBadge(mapInstanceRef.current, plotAssignments);
    } else {
      updateBadgeStyles(mapInstanceRef.current, playerRef.current);
    }

    map.on("click", function (evt) {
      let nextPrio = getLowestFreePriority(playerRef.current!);
      map.forEachFeatureAtPixel(evt.pixel, function (feature, layer) {
        if (
          feature &&
          feature.getGeometry()?.getType() === "Point" &&
          !lockedRef.current
        ) {
          const plot = feature.get("plot") as PlotData;
          if (!targetedRef.current) {
            updatePlayerPlot(plot.id, nextPrio);
            nextPrio++;
          } else {
            handleOpenModal(plot.id);
          }
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
          clearPlayerMappings={clearPlayerMappings}
          updatePlotAssignments={setPlotAssignments}
          targetedMode={targetedMode}
          setTargetedMode={setTargetedMode}
          contextDirty={contextDirty}
        />
        {/* Map Container */}
        <div ref={mapRef} id="map" className="map-container-style" />
        {mapReady && mapInstanceRef.current && (
          <MapHoverPopup map={mapInstanceRef.current} />
        )}
      </div>
      {playerRef.current && (
        <TargetedModal
          isOpen={isTargetedModalOpen}
          onClose={() => setIsTargetedModalOpen(false)}
          onSubmit={handleModalSubmit}
          plot={selectedPlot!}
          player={playerRef.current}
        />
      )}
    </div>
  );
}
