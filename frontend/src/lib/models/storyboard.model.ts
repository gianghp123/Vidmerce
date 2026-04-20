import type { ICampaign } from "./campaign.model";
import type { ISlide } from "./slide.model";

export interface IStoryboard {
  id: string;
  /**
   * ICampaign body copy
   */
  bodyCopy: string;
  /**
   * Call to action text
   */
  cta: string;
  /**
   * ICampaign headline
   */
  headline: string;
  /**
   * Associated job ID
   */
  jobId: string;
  /**
   * IStoryboard slides
   */
  slides: ISlide[];
}
