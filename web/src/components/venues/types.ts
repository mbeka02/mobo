/* ─── Venue Types ─── */
export interface VenueResponse {
  id: number;
  name: string;
  address: string;
  city: string;
  total_seats: number;
  created_at: string;
  updated_at?: string;
}

export interface VenuesAPIResponse {
  status: number;
  message: string;
  data: VenueResponse[] | null;
}

export interface CreateVenueAPIResponse {
  status: number;
  message: string;
  data: VenueResponse;
}

/* ─── Visual Helpers ─── */
const VENUE_GRADIENTS = [
  "linear-gradient(135deg, #AB3600 0%, #FF5F1F 100%)",
  "linear-gradient(135deg, #5C1A00 0%, #9B4425 100%)",
  "linear-gradient(135deg, #004B70 0%, #006493 100%)",
  "linear-gradient(135deg, #7C2E0F 0%, #FE916B 100%)",
  "linear-gradient(135deg, #003350 0%, #8DCDFF 100%)",
];

export function getVenueGradient(index: number) {
  return VENUE_GRADIENTS[index % VENUE_GRADIENTS.length];
}

export function getInitials(name: string) {
  return name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}
