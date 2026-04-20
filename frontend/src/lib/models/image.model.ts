import type { ImageStatus } from "../enums/image-status.enum";

export interface IImage {
  id: string;
  /**
   * S3 file key
   */
  fileKey: string;
  /**
   * IImage order
   */
  order: number;
  status: ImageStatus;
}
