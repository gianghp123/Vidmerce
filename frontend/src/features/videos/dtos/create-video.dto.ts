export interface VideoItemDto {
  assetId: string;
  duration: number;
  transition: string;
}

export interface CreateVideoDto {
  title: string;
  style: string;
  items: VideoItemDto[];
}