import { useState } from "react";
import { Play, Volume2, Captions, Maximize } from "lucide-react";
import type { Video } from "../models/video.model";
import { Hotspot } from "./Hotspot";

interface VideoPlayerProps {
  video: Video;
}

export function VideoPlayer({ video }: VideoPlayerProps) {
  const [isPlaying, setIsPlaying] = useState(false);
  const currentTime = 134;
  const duration = 290;

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  return (
    <div className="col-span-12 lg:col-span-8 group relative rounded-2xl overflow-hidden bg-surface-container-low shadow-2xl aspect-video">
      <img
        src={video.thumbnail_url}
        alt={video.title}
        className="w-full h-full object-cover brightness-90 group-hover:brightness-75 transition-all duration-700"
      />
      
      {video.hotspots.map((hotspot, index) => (
        <Hotspot key={index} hotspot={hotspot} />
      ))}
      
      <div className="absolute bottom-0 left-0 right-0 p-8 bg-gradient-to-t from-black/60 to-transparent flex items-center justify-between">
        <div className="flex items-center gap-6 text-white">
          <button
            onClick={() => setIsPlaying(!isPlaying)}
            className="cursor-pointer hover:opacity-80 transition-opacity"
          >
            <Play className="w-10 h-10 fill-current" />
          </button>
          <button className="cursor-pointer hover:opacity-80 transition-opacity">
            <Volume2 className="w-8 h-8" />
          </button>
          <span className="text-sm tracking-widest opacity-80">
            {formatTime(currentTime)} / {formatTime(duration)}
          </span>
        </div>
        <div className="flex items-center gap-4 text-white">
          <button className="cursor-pointer hover:opacity-80 transition-opacity">
            <Captions className="w-6 h-6" />
          </button>
          <button className="cursor-pointer hover:opacity-80 transition-opacity">
            <Maximize className="w-6 h-6" />
          </button>
        </div>
      </div>
    </div>
  );
}
