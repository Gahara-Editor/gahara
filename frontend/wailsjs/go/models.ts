export namespace audio {
	
	export class AudioNode {
	    type: string;
	    start: number;
	    end: number;
	    rid: string;
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new AudioNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.start = source["start"];
	        this.end = source["end"];
	        this.rid = source["rid"];
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

export namespace main {
	
	export class Video {
	    id: string;
	    name: string;
	    extension: string;
	    filepath: string;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new Video(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.extension = source["extension"];
	        this.filepath = source["filepath"];
	        this.duration = source["duration"];
	    }
	}

}

export namespace placeholder {
	
	export class PlaceholderNode {
	    type: string;
	    start: number;
	    end: number;
	    rid: string;
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new PlaceholderNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.start = source["start"];
	        this.end = source["end"];
	        this.rid = source["rid"];
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

export namespace timeline {
	
	export class Timeline {
	    timeline: any[][];
	
	    static createFrom(source: any = {}) {
	        return new Timeline(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.timeline = source["timeline"];
	    }
	}

}

export namespace video {
	
	export class ProcessingOpts {
	    resolution: string;
	    codec: string;
	    crf: string;
	    preset: string;
	    input_path?: string;
	    output_path: string;
	    filename: string;
	    video_format: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcessingOpts(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.resolution = source["resolution"];
	        this.codec = source["codec"];
	        this.crf = source["crf"];
	        this.preset = source["preset"];
	        this.input_path = source["input_path"];
	        this.output_path = source["output_path"];
	        this.filename = source["filename"];
	        this.video_format = source["video_format"];
	    }
	}
	export class VideoNode {
	    type: string;
	    start: number;
	    end: number;
	    rid: string;
	    id: string;
	    name: string;
	    losslessexport: boolean;
	
	    static createFrom(source: any = {}) {
	        return new VideoNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.start = source["start"];
	        this.end = source["end"];
	        this.rid = source["rid"];
	        this.id = source["id"];
	        this.name = source["name"];
	        this.losslessexport = source["losslessexport"];
	    }
	}

}

