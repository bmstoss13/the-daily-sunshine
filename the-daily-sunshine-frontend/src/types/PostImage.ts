export interface PostImage {
    id: string;
    post_id: string;
    image_url: string;
    image_description?: string;
    alt_text?: string;
    is_cover_image: boolean;
    created_at: Date;
}