export namespace backend {
	
	export class AudioDeviceInfo {
	    id: string;
	    name: string;
	    description: string;
	    isDefault: boolean;
	    isMonitor: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AudioDeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.isDefault = source["isDefault"];
	        this.isMonitor = source["isMonitor"];
	    }
	}
	export class EncoderInfo {
	    type: string;
	    name: string;
	    description: string;
	    devicePath?: string;
	    isHardware: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EncoderInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.type = source["type"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.devicePath = source["devicePath"];
	        this.isHardware = source["isHardware"];
	    }
	}
	export class RecordConfig {
	    resolution: string;
	    fps: number;
	    audioSource: string;
	    outputDir: string;
	    micDevice: string;
	    sysDevice: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.resolution = source["resolution"];
	        this.fps = source["fps"];
	        this.audioSource = source["audioSource"];
	        this.outputDir = source["outputDir"];
	        this.micDevice = source["micDevice"];
	        this.sysDevice = source["sysDevice"];
	    }
	}
	export class RecordedFileEntry {
	    filename: string;
	    filePath: string;
	    sizeHuman: string;
	    durationSec: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new RecordedFileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.filePath = source["filePath"];
	        this.sizeHuman = source["sizeHuman"];
	        this.durationSec = source["durationSec"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class RecordingStatus {
	    state: string;
	    durationSecs: number;
	    formatted: string;
	    outputFile?: string;
	    tempFile?: string;
	    encoderUsed?: string;
	    errorMessage?: string;
	
	    static createFrom(source: any = {}) {
	        return new RecordingStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.state = source["state"];
	        this.durationSecs = source["durationSecs"];
	        this.formatted = source["formatted"];
	        this.outputFile = source["outputFile"];
	        this.tempFile = source["tempFile"];
	        this.encoderUsed = source["encoderUsed"];
	        this.errorMessage = source["errorMessage"];
	    }
	}

}

