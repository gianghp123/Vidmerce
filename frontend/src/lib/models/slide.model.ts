import type { ImageRole } from "../enums/image-role.enum";

export interface ISlide {
  imageRole: ImageRole;
  /**
   * Overlay text on image
   */
  overlayText?: string;
  /**
   * S3 key for image
   */
  s3Key: string;
  /**
   * ISlide number
   */
  slideNumber: number;
}
