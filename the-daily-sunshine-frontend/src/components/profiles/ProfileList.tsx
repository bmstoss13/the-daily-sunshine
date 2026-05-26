import { useState } from "react";
import { useListProfiles } from "../../hooks/useProfiles";
import type { Profile } from "../../types/Profile";

export default function ProfileList(){
    const [offset, setOffset] = useState<number>(0)
    const { data, isPending } = useListProfiles(20, offset)
    
    return (
        <div>
            {data !== null && !isPending && (
                <div>
                    {Array.isArray(data) ? data.map((profile: Profile) => {
                        return (
                            <div key={profile.id}>
                                {profile.username}
                            </div>
                        )
                    }) : null}
                </div>
            )}
            <button
                onClick={() => setOffset(offset + 20)}
            >
                Don't press    
            </button>
        </div>
    )
}