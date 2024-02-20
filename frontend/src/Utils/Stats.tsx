import { GPUStats } from "../Data";

export function isFree(s: GPUStats): Boolean {
  return s.gpu_util < 5;
}
