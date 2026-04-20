import type { CampaignStatus } from "../enums/campaign-status.enum";
import type { IStoryboard } from "./storyboard.model";

export interface ICampaign {
  id: string;
  /**
   * Associated asset ID
   */
  assetId: string;
  /**
   * Creation timestamp
   */
  createdAt: string;
  status: CampaignStatus;
  /**
   * ICampaign storyboard
   */
  storyboard?: IStoryboard;
  /**
   * Last update timestamp
   */
  updatedAt?: string;
}
