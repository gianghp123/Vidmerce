import { JobStatus } from "../enums/job-status.enum";
import { JobType } from "../enums/job-type.enum";

export interface IJob {
  id: string;
  /**
   * Creation timestamp
   */
  createdAt: string;
  /**
   * Error log if failed
   */
  errorLog?: string;
  /**
   * IJob payload data
   */
  payload?: { [key: string]: any };
  status: JobStatus;
  /**
   * Target entity ID
   */
  targetId?: string;
  type: JobType;
}
