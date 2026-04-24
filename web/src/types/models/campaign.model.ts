import type { IStoryboard } from "./storyboard.model";
import { CampaignStatus } from "../enums/campaign-status.enum";

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
