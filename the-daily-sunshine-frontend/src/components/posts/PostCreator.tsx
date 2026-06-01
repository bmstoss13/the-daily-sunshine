import RichTextEditor from "../editor/RichTextEditor";

interface CreatePostProps {
    content: string;
    onChange: (content: string) => void;
}

export default function PostCreator({content, onChange}: CreatePostProps){
    return (
        <RichTextEditor 
            content={content}
            onChange={onChange}
        />
    )
}