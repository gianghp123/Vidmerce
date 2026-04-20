import type { AssetStatus } from "../enums/asset-status.enum";

export interface IAsset {
  id: string;
  /**
   * Creation timestamp
   */
  createdAt: string;
  /**
   * Number of images
   */
  imageCount: number;
  /**
   * Product name
   */
  name: string;
  /**
   * Product price
   */
  price: number;
  /**
   * Product URL
   */
  productUrl: string;
  status: AssetStatus;
}
